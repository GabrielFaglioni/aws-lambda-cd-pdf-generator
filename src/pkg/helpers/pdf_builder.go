package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/FavorDespaches/content-declaration-pdf-generator/pkg/types"
	"github.com/jung-kurt/gofpdf"
)

const (
	paddingBetweenDC = 1.0
	nObjetosPorFolha = 4
)

func GenerateContentDeclarationPDFLocal(solicitarDeclaracaoConteudo types.SolicitarDeclaracaoConteudo) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)
	objetosPostais, objErr := rearrangeObjetoPostal(solicitarDeclaracaoConteudo.ObjetosPostais)

	if objErr != nil {
		fmt.Println("Error rearranging ObjetoPostal, using original: ", objErr)
		objetosPostais = solicitarDeclaracaoConteudo.ObjetosPostais
	}

	nObjetosPostais := len(objetosPostais)
	remetente := solicitarDeclaracaoConteudo.Remetente
	skipNext := false

	fmt.Println(" - Número de Declarações de Conteúdo: ", nObjetosPostais)
	for i, objetoPostal := range objetosPostais {
		if skipNext {
			// Skip this iteration and reset the flag
			skipNext = false
			continue
		}

		fmt.Println("   - Desenhando a D.C. ", i)
		pdf.AddPage()

		x := 0.0
		y := 0.0
		nextY := DrawContentDeclarationV2(pdf, remetente, objetoPostal, true, x, y)
		nObjetos := len(objetoPostal.DeclaracoesConteudo)

		fmt.Printf("   - Para %d Objetos:\n", nObjetos)
		fmt.Printf("     - nextY: %.2f\n", nextY)

		if i < nObjetosPostais-1 {
			nObjetosProxDC := len(objetosPostais[i+1].DeclaracoesConteudo)

			if nObjetosProxDC+nObjetos <= nObjetosPorFolha {
				fmt.Println("   - Desenhando a D.C. ", i+1)
				nextY += paddingBetweenDC
				DrawDashedLine(pdf, nextY+paddingBetweenDC)
				nextY = DrawContentDeclarationV2(pdf, remetente, objetosPostais[i+1], true, x, nextY)
				fmt.Printf("   - Para %d Objetos:\n", nObjetosProxDC)
				fmt.Printf("     - nextY: %.2f\n", nextY)
				skipNext = true
			}
		}
	}

	return pdf.OutputFileAndClose("label.pdf")
}

func GenerateContentDeclaration(solicitarDeclaracaoConteudo types.SolicitarDeclaracaoConteudo) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)
	objetosPostais, objErr := rearrangeObjetoPostal(solicitarDeclaracaoConteudo.ObjetosPostais)

	if objErr != nil {
		fmt.Println("Error rearranging ObjetoPostal, using original: ", objErr)
		objetosPostais = solicitarDeclaracaoConteudo.ObjetosPostais
	}

	nObjetosPostais := len(objetosPostais)
	remetente := solicitarDeclaracaoConteudo.Remetente
	skipNext := false

	fmt.Println(" - Número de Páginas do PDF: ", len(objetosPostais))
	for i, objetoPostal := range objetosPostais {
		if skipNext {
			// Skip this iteration and reset the flag
			skipNext = false
			continue
		}

		fmt.Println("   - Desenhando a D.C. ", i)
		pdf.AddPage()

		x := 0.0
		y := 0.0
		nextY := DrawContentDeclarationV2(pdf, remetente, objetoPostal, false, x, y)
		nObjetos := len(objetoPostal.DeclaracoesConteudo)

		fmt.Printf("   - Para %d Objetos:\n", nObjetos)
		fmt.Printf("     - nextY: %.2f\n", nextY)

		if i < nObjetosPostais-1 {
			nObjetosProxDC := len(objetosPostais[i+1].DeclaracoesConteudo)

			if nObjetosProxDC+nObjetos <= nObjetosPorFolha {
				fmt.Println("   - Desenhando a D.C. ", i+1)
				nextY += paddingBetweenDC
				DrawDashedLine(pdf, nextY+paddingBetweenDC)
				nextY = DrawContentDeclarationV2(pdf, remetente, objetosPostais[i+1], true, x, nextY)
				fmt.Printf("   - Para %d Objetos:\n", nObjetosProxDC)
				fmt.Printf("     - nextY: %.2f\n", nextY)
				skipNext = true
			}
		}
	}

	var buffer bytes.Buffer

	err := pdf.Output(&buffer)
	if err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return base64Str, nil
}

type PackageWeight struct {
	Signal string // "g" | "Kg"
	Value  float64
}

func DrawContentDeclarationV2(pdf *gofpdf.Fpdf, remetente types.SolicitarEtiquetaRemetente, objetoPostal types.SolicitarDeclaracaoConteudoObjetoPostal, local bool, x, y float64) float64 {
	destinatario := objetoPostal.Destinatario
	declaracoesConteudo := objetoPostal.DeclaracoesConteudo

	packageWeight := PackageWeight{
		Signal: "g",
		Value:  objetoPostal.Peso,
	}

	if packageWeight.Value > 1000 {
		packageWeight.Signal = "Kg"
		packageWeight.Value = packageWeight.Value / 1000.0
	}

	nextY := DrawHeader(pdf, x, y)
	nextY = DrawSenderReceiverDataV2(pdf, x, nextY, remetente, destinatario)
	nextY = DrawDeclarationItemsV2(pdf, x, nextY, declaracoesConteudo, packageWeight)
	nextY = DrawDeclaration(pdf, x, nextY, remetente.CidadeRemetente)
	nextY = DrawObservations(pdf, x, nextY)
	return nextY
}

func rearrangeObjetoPostal(objetos []types.SolicitarDeclaracaoConteudoObjetoPostal) (rearranged []types.SolicitarDeclaracaoConteudoObjetoPostal, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert the panic to an error
			err = fmt.Errorf("rearrangeObjetoPostalV2 panicked: %v", r)
		}
	}()

	// Step 1: Sort the ObjetoPostal based on the length of DeclaracaoConteudo.
	sort.Slice(objetos, func(i, j int) bool {
		return len(objetos[i].DeclaracoesConteudo) < len(objetos[j].DeclaracoesConteudo)
	})

	// Step 2: Create buckets for different lengths of DeclaracaoConteudo.
	buckets := make(map[int][]types.SolicitarDeclaracaoConteudoObjetoPostal)
	for _, obj := range objetos {
		length := len(obj.DeclaracoesConteudo)
		buckets[length] = append(buckets[length], obj)
	}

	// Step 3: Merge the objects from the buckets, attempting to pair objects with similar lengths.
	var result []types.SolicitarDeclaracaoConteudoObjetoPostal
	for i := 1; i <= 4; i++ {
		for len(buckets[i]) > 0 {
			// Add the object from the current bucket.
			result = append(result, buckets[i][0])
			buckets[i] = buckets[i][1:]

			// Determine the length needed to pair with the current object.
			neededLength := 4 - i
			if neededLength > 0 && len(buckets[neededLength]) > 0 {
				// If there's an object with the needed length, pair it with the current object.
				result = append(result, buckets[neededLength][0])
				buckets[neededLength] = buckets[neededLength][1:]
			}
		}
	}

	// Step 4: Add any remaining objects (this handles cases where an exact pairing wasn't possible).
	for _, remainingObjects := range buckets {
		result = append(result, remainingObjects...)
	}

	return result, nil
}
