package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/FavorDespaches/content-declaration-pdf-generator/pkg/types"
	"github.com/jung-kurt/gofpdf"
)

const (
	pageWidth              = 210.0 // A4 width in mm
	pageHeight             = 297.0 // A4 height in mm
	defaultBorderLineWidth = 0.6
	defaultLineWidth       = 0.2
	defaultPagePadding     = 3.0
	defaultGridPadding     = defaultPagePadding / 2
	declarationBlockHeight = 48.0
	observationBlockHeight = 13.0
)

func DrawDashedLine(pdf *gofpdf.Fpdf, y float64) {
	pdf.SetDashPattern([]float64{3, 3}, 0)
	pdf.SetLineWidth(0.1)
	pdf.Line(0, y, pageWidth, y)
	pdf.SetDashPattern([]float64{}, 0)
	pdf.SetLineWidth(0.3)
}

// ! ==================== CONTENT DECLARATION HEADER ====================
func DrawHeader(pdf *gofpdf.Fpdf, x, y float64) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	xHeader := x + defaultPagePadding
	yHeader := y + defaultPagePadding

	fontSize := 12.0
	headerHeight := 1.6 * pdf.PointConvert(fontSize)

	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(xHeader, yHeader, pageWidth-2*xHeader, headerHeight, "D")
	pdf.SetLineWidth(defaultLineWidth)

	pdf.SetFont("Arial", "B", fontSize)
	headerString := "DECLARAÇÃO DE CONTEÚDO"

	headerTextX := (pageWidth - 2*defaultPagePadding - pdf.GetStringWidth(headerString)) / 2
	headerTextY := yHeader + headerHeight/2 + pdf.PointConvert(fontSize)/2 - 0.6

	pdf.Text(headerTextX, headerTextY, translator(headerString))

	nextY := headerHeight + yHeader

	return nextY
}

// ! ==================== SENDER / RECEIVER ADDRESS DATA ====================
func getEnderecoString(remetente types.Remetente, destinatario types.Destinatario, nacional types.Nacional) types.SenderReceiverAddressData {
	//- CONSTRUINDO NOME
	senderName := remetente.NomeRemetente
	receiverName := destinatario.NomeDestinatario

	//- CONSTRUINDO ENDERECO
	senderAddress := remetente.LogradouroRemetente
	if remetente.NumeroRemetente != "" {
		senderAddress += ", "
		senderAddress += remetente.NumeroRemetente
	}

	receiverAddress := destinatario.LogradouroDestinatario
	if destinatario.NumeroEndDestinatario != "" {
		receiverAddress += ", "
		receiverAddress += destinatario.NumeroEndDestinatario
	}

	//- CONSTRUINDO COMPLEMENTO / BAIRRO
	var senderComplementoBairro string
	if remetente.ComplementoRemetente != "" {
		senderComplementoBairro += remetente.ComplementoRemetente
	}
	if remetente.ComplementoRemetente != "" && remetente.BairroRemetente != "" {
		senderComplementoBairro += " / "
	}
	if remetente.BairroRemetente != "" {
		senderComplementoBairro += remetente.BairroRemetente
	}

	var receiverComplementoBairro string
	if destinatario.ComplementoDestinatario != "" {
		receiverComplementoBairro += destinatario.ComplementoDestinatario
	}
	if destinatario.ComplementoDestinatario != "" && nacional.BairroDestinatario != "" {
		receiverComplementoBairro += " / "
	}
	if nacional.BairroDestinatario != "" {
		receiverComplementoBairro += nacional.BairroDestinatario
	}

	//- CONSTRUINDO CPFCNPJ
	var senderCPFCNPJ string
	if remetente.CpfCnpjRemetente == 0 {
		senderCPFCNPJ = "-"
	} else {
		senderCPFCNPJ = formatCPFCNPJ(remetente.CpfCnpjRemetente)
	}

	var receiverCPFCNPJ string
	if destinatario.CpfCnpjDestinatario == 0 {
		receiverCPFCNPJ = "-"
	} else {
		receiverCPFCNPJ = formatCPFCNPJ(destinatario.CpfCnpjDestinatario)
	}

	//- CONSTRUINDO ADDRESSDATA
	senderAddressData := types.AddressData{
		Nome:              senderName,
		Endereco:          senderAddress,
		ComplementoBairro: senderComplementoBairro,
		Uf:                remetente.UfRemetente,
		Cidade:            remetente.CidadeRemetente,
		Cep:               remetente.CepRemetente,
		CPFCNPJ:           senderCPFCNPJ,
	}

	receiverAddressData := types.AddressData{
		Nome:              receiverName,
		Endereco:          receiverAddress,
		ComplementoBairro: receiverComplementoBairro,
		Uf:                nacional.UfDestinatario,
		Cidade:            nacional.CidadeDestinatario,
		Cep:               nacional.CepDestinatario,
		CPFCNPJ:           receiverCPFCNPJ,
	}

	return types.SenderReceiverAddressData{
		SenderAddressData:   senderAddressData,
		ReceiverAddressData: receiverAddressData,
	}
}

func getEnderecoStringV2(remetente types.SolicitarEtiquetaRemetente, destinatario types.SolicitarEtiquetaDestinatario) types.SenderReceiverAddressData {
	//- CONSTRUINDO NOME
	senderName := remetente.NomeRemetente
	receiverName := destinatario.NomeDestinatario

	//- CONSTRUINDO ENDERECO
	senderAddress := remetente.LogradouroRemetente
	if remetente.NumeroRemetente != "" {
		senderAddress += ", "
		senderAddress += remetente.NumeroRemetente
	}

	receiverAddress := destinatario.LogradouroDestinatario
	if destinatario.NumeroDestinatario != "" {
		receiverAddress += ", "
		receiverAddress += destinatario.NumeroDestinatario
	}

	//- CONSTRUINDO COMPLEMENTO / BAIRRO
	var senderComplementoBairro string
	if remetente.ComplementoRemetente != nil {
		senderComplementoBairro += *remetente.ComplementoRemetente
	}
	if remetente.ComplementoRemetente != nil && remetente.BairroRemetente != "" {
		senderComplementoBairro += " / "
	}
	if remetente.BairroRemetente != "" {
		senderComplementoBairro += remetente.BairroRemetente
	}

	var receiverComplementoBairro string
	if destinatario.ComplementoDestinatario != nil {
		receiverComplementoBairro += *destinatario.ComplementoDestinatario
	}
	if destinatario.ComplementoDestinatario != nil && destinatario.BairroDestinatario != "" {
		receiverComplementoBairro += " / "
	}
	if destinatario.BairroDestinatario != "" {
		receiverComplementoBairro += destinatario.BairroDestinatario
	}

	//- CONSTRUINDO CPFCNPJ
	var senderCPFCNPJ string
	if remetente.CpfCnpjRemetente == "" {
		senderCPFCNPJ = "-"
	} else {
		senderCPFCNPJ = formatCPFCNPJV2(remetente.CpfCnpjRemetente)
	}

	var receiverCPFCNPJ string
	if *destinatario.CpfCnpjDestinatario == "" {
		receiverCPFCNPJ = "-"
	} else {
		receiverCPFCNPJ = formatCPFCNPJV2(*destinatario.CpfCnpjDestinatario)
	}

	//- CONSTRUINDO ADDRESSDATA
	senderAddressData := types.AddressData{
		Nome:              senderName,
		Endereco:          senderAddress,
		ComplementoBairro: senderComplementoBairro,
		Uf:                remetente.UfRemetente,
		Cidade:            remetente.CidadeRemetente,
		Cep:               remetente.CepRemetente,
		CPFCNPJ:           senderCPFCNPJ,
	}

	receiverAddressData := types.AddressData{
		Nome:              receiverName,
		Endereco:          receiverAddress,
		ComplementoBairro: receiverComplementoBairro,
		Uf:                destinatario.UfDestinatario,
		Cidade:            destinatario.CidadeDestinatario,
		Cep:               destinatario.CepDestinatario,
		CPFCNPJ:           receiverCPFCNPJ,
	}

	return types.SenderReceiverAddressData{
		SenderAddressData:   senderAddressData,
		ReceiverAddressData: receiverAddressData,
	}
}

func DrawSenderReceiverData(pdf *gofpdf.Fpdf, x, y float64, remetente types.Remetente, destinatario types.Destinatario, nacional types.Nacional) float64 {
	senderReceiverBlockWidth := (pageWidth - 2*defaultPagePadding - defaultGridPadding) / 2
	senderReceiverBlockHeight := 35.0

	senderX := x + defaultPagePadding
	senderY := y + defaultGridPadding

	senderReceiverAddressData := getEnderecoString(remetente, destinatario, nacional)
	senderAddressData := senderReceiverAddressData.SenderAddressData
	receiverAddressData := senderReceiverAddressData.ReceiverAddressData

	DrawSenderReceiverBlock(pdf, senderX, senderY, senderReceiverBlockWidth, senderReceiverBlockHeight, "REMETENTE", senderAddressData)

	receiverX := x + defaultPagePadding + senderReceiverBlockWidth + defaultGridPadding
	receiverY := y + defaultGridPadding

	DrawSenderReceiverBlock(pdf, receiverX, receiverY, senderReceiverBlockWidth, senderReceiverBlockHeight, "DESTINATÁRIO", receiverAddressData)
	nextY := y + senderReceiverBlockHeight
	return nextY
}

func DrawSenderReceiverDataV2(pdf *gofpdf.Fpdf, x, y float64, remetente types.SolicitarEtiquetaRemetente, destinatario types.SolicitarEtiquetaDestinatario) float64 {
	senderReceiverBlockWidth := (pageWidth - 2*defaultPagePadding - defaultGridPadding) / 2
	senderReceiverBlockHeight := 35.0

	senderX := x + defaultPagePadding
	senderY := y + defaultGridPadding

	senderReceiverAddressData := getEnderecoStringV2(remetente, destinatario)
	senderAddressData := senderReceiverAddressData.SenderAddressData
	receiverAddressData := senderReceiverAddressData.ReceiverAddressData

	DrawSenderReceiverBlock(pdf, senderX, senderY, senderReceiverBlockWidth, senderReceiverBlockHeight, "REMETENTE", senderAddressData)

	receiverX := x + defaultPagePadding + senderReceiverBlockWidth + defaultGridPadding
	receiverY := y + defaultGridPadding

	DrawSenderReceiverBlock(pdf, receiverX, receiverY, senderReceiverBlockWidth, senderReceiverBlockHeight, "DESTINATÁRIO", receiverAddressData)
	nextY := y + senderReceiverBlockHeight
	return nextY
}

func DrawSenderReceiverBlock(pdf *gofpdf.Fpdf, x, y, w, h float64, person string, addressData types.AddressData) float64 {
	const firstRowField = "NOME: "
	const secondRowField = "ENDEREÇO:"
	const thirdRowField = "COMP./BAIRRO: "

	const fourthRowFirstColumnField = "UF:  "
	const fourthRowSecondColumnField = "CIDADE:    "

	const fifthRowFirstColumnField = "CEP:  "
	const fifthRowSecondColumnField = "CPF/CNPJ:    "

	translator := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(x, y, w, h, "D")
	pdf.SetLineWidth(defaultLineWidth)

	//! DRAW BLOCK HEADER
	headerTextSize := 10.0
	pdf.SetFont("Arial", "B", headerTextSize)

	headerX := x + (w-pdf.GetStringWidth(person))/2
	headerY := y + pdf.PointConvert(headerTextSize) + defaultGridPadding/2

	pdf.Text(headerX, headerY, translator(person))

	headerLineY := headerY + defaultGridPadding/2 + 1.0
	pdf.Line(x, headerLineY, x+w, headerLineY)

	//! DRAW LINES
	lineX := x + defaultGridPadding/2
	lineHeight := (h - (headerY - y + defaultGridPadding)) / 5

	nextY := DrawSenderReceiverBlockSingleColumnLine(pdf, firstRowField, addressData.Nome, lineX, headerLineY, w, lineHeight)
	nextY = DrawSenderReceiverBlockSingleColumnLine(pdf, secondRowField, addressData.Endereco, lineX, nextY, w, lineHeight)
	nextY = DrawSenderReceiverBlockSingleColumnLine(pdf, thirdRowField, addressData.ComplementoBairro, lineX, nextY, w, lineHeight)
	nextY = DrawSenderReceiverBlockDoubleColumnLine(pdf, fourthRowFirstColumnField, addressData.Uf, fourthRowSecondColumnField, addressData.Cidade, lineX, nextY, w, lineHeight)
	nextY = DrawSenderReceiverBlockDoubleColumnLine(pdf, fifthRowFirstColumnField, formatCEP(addressData.Cep), fifthRowSecondColumnField, addressData.CPFCNPJ, lineX, nextY, w, lineHeight)

	return nextY
}

func DrawSenderReceiverBlockSingleColumnLine(pdf *gofpdf.Fpdf, field string, value string, x, y, w, h float64) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	fieldTextSize := 9.0
	pdf.SetFont("Arial", "B", fieldTextSize)

	lineY := y + pdf.PointConvert(fieldTextSize) + defaultGridPadding/2
	pdf.Text(x, lineY, translator(field))

	valueX := x + pdf.GetStringWidth(field)
	valueTextSize := 9.0

	if len(value) > 50 {
		valueTextSize = 7.0
	}

	pdf.SetFont("Arial", "", valueTextSize)
	pdf.Text(valueX, lineY, translator(value))

	nextY := y + h
	pdf.Line(x-defaultGridPadding/2, nextY, x+w-defaultGridPadding/2, nextY)
	return nextY
}

func DrawSenderReceiverBlockDoubleColumnLine(pdf *gofpdf.Fpdf, field1 string, value1 string, field2 string, value2 string, x, y, w, h float64) float64 {
	//! FIRST COLUMN
	firstColumnX := x - defaultGridPadding/2
	firstColumnWidth := 3 * w / 12.0
	pdf.Rect(firstColumnX, y, firstColumnWidth, h, "D")

	pdf.SetFont("Arial", "B", 9.0)
	drawYCenteredText(pdf, x, y, h, 9.0, field1)

	pdf.SetFont("Arial", "", 9.0)
	firstColumnValueX := x + pdf.GetStringWidth(field1) - defaultGridPadding/2
	drawYCenteredText(pdf, firstColumnValueX, y, h, 9.0, value1)

	//! SECOND COLUMN
	secondColumnX := firstColumnX + firstColumnWidth
	secondColumnWidth := w - firstColumnWidth
	pdf.Rect(secondColumnX, y, secondColumnWidth, h, "D")
	secondColumnX += defaultGridPadding / 2

	pdf.SetFont("Arial", "B", 9.0)
	drawYCenteredText(pdf, secondColumnX, y, h, 9.0, field2)

	pdf.SetFont("Arial", "", 8.0)
	secondColumnValueX := secondColumnX + pdf.GetStringWidth(field2) - defaultGridPadding/2
	drawYCenteredText(pdf, secondColumnValueX, y, h, 9.0, value2)

	nextY := y + h
	return nextY
}

// ! ==================== DECLARATION ITEMS ====================
func DrawDeclarationItems(pdf *gofpdf.Fpdf, x, y float64, declaracaoConteudo []types.DeclaracaoConteudo, packageWeight PackageWeight) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	declarationItemsHeaderTextSize := 11.0
	declarationItemsHeaderHeight := 1.6 * pdf.PointConvert(declarationItemsHeaderTextSize)
	declarationItemsRowHeight := 1.6 * pdf.PointConvert(10.0)

	n := float64(len(declaracaoConteudo))
	//- ========== DECLARATION ITEMS BOX ==========
	declarationItemsX := x + defaultPagePadding
	declarationItemsY := y + 2*defaultGridPadding
	declarationItemsWidth := pageWidth - 2*defaultPagePadding
	declarationItemsHeight := declarationItemsHeaderHeight + 3*declarationItemsRowHeight + n*declarationItemsRowHeight

	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(declarationItemsX, declarationItemsY, declarationItemsWidth, declarationItemsHeight, "D")

	//- ========== DECLARATION ITEMS HEADER ==========
	declarationItemsHeaderString := "IDENTIFICAÇÃO DOS BENS"

	pdf.SetFont("Arial", "B", declarationItemsHeaderTextSize)

	declarationItemsHeaderPadding := 1.2 * pdf.PointConvert(declarationItemsHeaderTextSize)
	declarationItemsHeaderTextX := (pageWidth - 2*defaultPagePadding - pdf.GetStringWidth(declarationItemsHeaderString)) / 2
	declarationItemsHeaderTextY := declarationItemsY + declarationItemsHeaderPadding
	pdf.Text(declarationItemsHeaderTextX, declarationItemsHeaderTextY, translator(declarationItemsHeaderString))
	nextY := declarationItemsHeaderTextY + 0.4*declarationItemsHeaderPadding
	pdf.Line(declarationItemsX, nextY, pageWidth-defaultPagePadding, nextY)
	pdf.SetLineWidth(defaultLineWidth)

	//- DECLARATION ITEMS COLUMNS
	nextY = DrawDeclarationItemsColumns(pdf, x, nextY, declarationItemsWidth)
	var totalPrice float64
	var totalQuantity int
	//- DECLARATION ITEMS ROWS
	for i, declaracaoConteudoItem := range declaracaoConteudo {
		nextY = DrawDeclarationItemsRow(pdf, x, nextY, declarationItemsWidth, declaracaoConteudoItem, i+1)
		totalPrice += declaracaoConteudoItem.ValorUnitario * float64(declaracaoConteudoItem.Quantidade)
		totalQuantity += declaracaoConteudoItem.Quantidade
	}

	//- DECLARATION ITEMS RESUME
	DrawDeclarationItemsResume(pdf, x, nextY, declarationItemsWidth, totalQuantity, totalPrice, packageWeight)
	return declarationItemsY + declarationItemsHeight
}

func DrawDeclarationItemsV2(pdf *gofpdf.Fpdf, x, y float64, declaracoesConteudo []types.SolicitarEtiquetaDeclaracaoConteudo, packageWeight PackageWeight) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	declarationItemsHeaderTextSize := 11.0
	declarationItemsHeaderHeight := 1.6 * pdf.PointConvert(declarationItemsHeaderTextSize)
	declarationItemsRowHeight := 1.6 * pdf.PointConvert(10.0)

	n := float64(len(declaracoesConteudo))
	//- ========== DECLARATION ITEMS BOX ==========
	declarationItemsX := x + defaultPagePadding
	declarationItemsY := y + 2*defaultGridPadding
	declarationItemsWidth := pageWidth - 2*defaultPagePadding
	declarationItemsHeight := declarationItemsHeaderHeight + 3*declarationItemsRowHeight + n*declarationItemsRowHeight

	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(declarationItemsX, declarationItemsY, declarationItemsWidth, declarationItemsHeight, "D")

	//- ========== DECLARATION ITEMS HEADER ==========
	declarationItemsHeaderString := "IDENTIFICAÇÃO DOS BENS"

	pdf.SetFont("Arial", "B", declarationItemsHeaderTextSize)

	declarationItemsHeaderPadding := 1.2 * pdf.PointConvert(declarationItemsHeaderTextSize)
	declarationItemsHeaderTextX := (pageWidth - 2*defaultPagePadding - pdf.GetStringWidth(declarationItemsHeaderString)) / 2
	declarationItemsHeaderTextY := declarationItemsY + declarationItemsHeaderPadding
	pdf.Text(declarationItemsHeaderTextX, declarationItemsHeaderTextY, translator(declarationItemsHeaderString))
	nextY := declarationItemsHeaderTextY + 0.4*declarationItemsHeaderPadding
	pdf.Line(declarationItemsX, nextY, pageWidth-defaultPagePadding, nextY)
	pdf.SetLineWidth(defaultLineWidth)

	//- DECLARATION ITEMS COLUMNS
	nextY = DrawDeclarationItemsColumns(pdf, x, nextY, declarationItemsWidth)
	var totalPrice float64
	var totalQuantity int
	//- DECLARATION ITEMS ROWS
	for i, declaracaoConteudoItem := range declaracoesConteudo {
		nextY = DrawDeclarationItemsRowV2(pdf, x, nextY, declarationItemsWidth, declaracaoConteudoItem, i+1)
		totalPrice += declaracaoConteudoItem.ValorUnitario * float64(declaracaoConteudoItem.Quantidade)
		totalQuantity += declaracaoConteudoItem.Quantidade
	}

	//- DECLARATION ITEMS RESUME
	DrawDeclarationItemsResume(pdf, x, nextY, declarationItemsWidth, totalQuantity, totalPrice, packageWeight)
	return declarationItemsY + declarationItemsHeight
}

func DrawDeclarationItemsColumns(pdf *gofpdf.Fpdf, x, y, w float64) float64 {
	//translator := pdf.UnicodeTranslatorFromDescriptor("")
	textSize := 10.0
	pdf.SetFont("Arial", "B", textSize)
	columnY := y
	columnHeight := 1.6 * pdf.PointConvert(textSize)

	column1X := x + defaultPagePadding
	column1Width := 1.5 * w / 12
	pdf.Rect(column1X, columnY, column1Width, columnHeight, "D")
	drawXYCenteredText(pdf, column1X, columnY, column1Width, columnHeight, textSize, "ITEM")

	column2X := column1X + column1Width
	column2Width := 6.5 * w / 12
	pdf.Rect(column2X, columnY, column2Width, columnHeight, "D")
	drawXYCenteredText(pdf, column2X, columnY, column2Width, columnHeight, textSize, "CONTEUDO")

	column3X := column2X + column2Width
	column3Width := 1.5 * w / 12
	pdf.Rect(column3X, columnY, column3Width, columnHeight, "D")
	drawXYCenteredText(pdf, column3X, columnY, column3Width, columnHeight, textSize, "QUANT.")

	column4X := column3X + column3Width
	column4Width := 2.5 * w / 12
	pdf.Rect(column4X, columnY, column4Width, columnHeight, "D")
	drawXYCenteredText(pdf, column4X, columnY, column4Width, columnHeight, textSize, "VALOR (R$)")

	return y + columnHeight
}

func DrawDeclarationItemsRow(pdf *gofpdf.Fpdf, x, y, w float64, declaracaoConteudoItem types.DeclaracaoConteudo, index int) float64 {
	textSize := 10.0
	pdf.SetFont("Arial", "", textSize)
	columnY := y
	columnHeight := 1.6 * pdf.PointConvert(textSize)

	column1X := x + defaultPagePadding
	column1Width := 1.5 * w / 12
	pdf.Rect(column1X, columnY, column1Width, columnHeight, "D")
	drawXYCenteredText(pdf, column1X, columnY, column1Width, columnHeight, textSize, fmt.Sprintf("%d", index))

	column2X := column1X + column1Width
	column2Width := 6.5 * w / 12
	pdf.Rect(column2X, columnY, column2Width, columnHeight, "D")
	drawXYCenteredText(pdf, column2X, columnY, column2Width, columnHeight, textSize, declaracaoConteudoItem.Conteudo)

	column3X := column2X + column2Width
	column3Width := 1.5 * w / 12
	pdf.Rect(column3X, columnY, column3Width, columnHeight, "D")
	drawXYCenteredText(pdf, column3X, columnY, column3Width, columnHeight, textSize, fmt.Sprintf("%d", declaracaoConteudoItem.Quantidade))

	column4X := column3X + column3Width
	column4Width := 2.5 * w / 12
	pdf.Rect(column4X, columnY, column4Width, columnHeight, "D")

	valorUnitarioStr := fmt.Sprintf("%.2f", declaracaoConteudoItem.ValorUnitario)
	valorUnitario := strings.Replace(valorUnitarioStr, ".", ",", -1)
	drawXYCenteredText(pdf, column4X, columnY, column4Width, columnHeight, textSize, valorUnitario)

	return y + columnHeight
}

func DrawDeclarationItemsRowV2(pdf *gofpdf.Fpdf, x, y, w float64, declaracaoConteudoItem types.SolicitarEtiquetaDeclaracaoConteudo, index int) float64 {
	textSize := 10.0
	pdf.SetFont("Arial", "", textSize)
	columnY := y
	columnHeight := 1.6 * pdf.PointConvert(textSize)

	column1X := x + defaultPagePadding
	column1Width := 1.5 * w / 12
	pdf.Rect(column1X, columnY, column1Width, columnHeight, "D")
	drawXYCenteredText(pdf, column1X, columnY, column1Width, columnHeight, textSize, fmt.Sprintf("%d", index))

	column2X := column1X + column1Width
	column2Width := 6.5 * w / 12
	pdf.Rect(column2X, columnY, column2Width, columnHeight, "D")
	drawXYCenteredText(pdf, column2X, columnY, column2Width, columnHeight, textSize, declaracaoConteudoItem.Conteudo)

	column3X := column2X + column2Width
	column3Width := 1.5 * w / 12
	pdf.Rect(column3X, columnY, column3Width, columnHeight, "D")
	drawXYCenteredText(pdf, column3X, columnY, column3Width, columnHeight, textSize, fmt.Sprintf("%d", declaracaoConteudoItem.Quantidade))

	column4X := column3X + column3Width
	column4Width := 2.5 * w / 12
	pdf.Rect(column4X, columnY, column4Width, columnHeight, "D")

	valorUnitarioStr := fmt.Sprintf("%.2f", declaracaoConteudoItem.ValorUnitario)
	valorUnitario := strings.Replace(valorUnitarioStr, ".", ",", -1)
	drawXYCenteredText(pdf, column4X, columnY, column4Width, columnHeight, textSize, valorUnitario)

	return y + columnHeight
}

func DrawDeclarationItemsResume(pdf *gofpdf.Fpdf, x, y, w float64, totalQuantity int, totalPrice float64, packageWeight PackageWeight) float64 {
	textSize := 10.0
	pdf.SetFont("Arial", "B", textSize)
	rowHeight := 1.6 * pdf.PointConvert(textSize)

	row1Y := y

	row1Column1X := x + defaultPagePadding
	row1Column1Width := 8 * w / 12
	pdf.Rect(row1Column1X, row1Y, row1Column1Width, rowHeight, "D")
	drawYCenteredXJustifiedRightText(pdf, row1Column1X, row1Y, row1Column1Width, rowHeight, textSize, "TOTAIS")

	pdf.SetFont("Arial", "", textSize)
	row1Column2X := row1Column1X + row1Column1Width
	row1Column2Width := 1.5 * w / 12
	pdf.Rect(row1Column2X, row1Y, row1Column2Width, rowHeight, "D")
	drawXYCenteredText(pdf, row1Column2X, row1Y, row1Column2Width, rowHeight, textSize, fmt.Sprintf("%d", totalQuantity))

	row1Column3X := row1Column2X + row1Column2Width
	row1Column3Width := 2.5 * w / 12
	pdf.Rect(row1Column3X, row1Y, row1Column3Width, rowHeight, "D")
	drawXYCenteredText(pdf, row1Column3X, row1Y, row1Column3Width, rowHeight, textSize, strings.Replace(fmt.Sprintf("%.2f", totalPrice), ".", ",", -1))

	row2Y := y + rowHeight

	pdf.SetFont("Arial", "B", textSize)
	row2Column1X := row1Column1X
	row2Column1Width := 8 * w / 12
	row2Column1String := fmt.Sprintf("PESO TOTAL (%s)", packageWeight.Signal)
	pdf.Rect(row2Column1X, row2Y, row2Column1Width, rowHeight, "D")
	drawYCenteredXJustifiedRightText(pdf, row2Column1X, row2Y, row2Column1Width, rowHeight, textSize, row2Column1String)

	pdf.SetFont("Arial", "", textSize)
	row2Column2X := row2Column1X + row2Column1Width
	row2Column2Width := 4 * w / 12
	row2Column2String := strings.Replace(fmt.Sprintf("%.3f", packageWeight.Value), ".", ",", -1)
	pdf.Rect(row2Column2X, row2Y, row2Column2Width, rowHeight, "D")
	drawXYCenteredText(pdf, row2Column2X, row2Y, row2Column2Width, rowHeight, textSize, row2Column2String)

	return y + 2*rowHeight
}

// ! ==================== DECLARATION BLOCK ====================
func DrawDeclaration(pdf *gofpdf.Fpdf, x, y float64, city string) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")

	//! ========== DECLARATION BOX ==========
	declarationX := x + defaultPagePadding
	declarationY := y + defaultGridPadding

	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(declarationX, declarationY, pageWidth-2*defaultPagePadding, declarationBlockHeight, "D")

	//! ========== DECLARATION HEADER ==========
	declarationHeaderString := "DECLARAÇÃO"
	declarationHeaderTextSize := 11.0
	pdf.SetFont("Arial", "B", declarationHeaderTextSize)

	declarationHeaderPadding := 1.2 * pdf.PointConvert(declarationHeaderTextSize)
	declarationHeaderX := (pageWidth - 2*defaultPagePadding - pdf.GetStringWidth(declarationHeaderString)) / 2
	declarationHeaderY := declarationY + declarationHeaderPadding
	pdf.Text(declarationHeaderX, declarationHeaderY, translator(declarationHeaderString))
	nextY := declarationHeaderY + 0.4*pdf.PointConvert(declarationHeaderTextSize)
	pdf.Line(declarationX, nextY, pageWidth-defaultPagePadding, nextY)
	pdf.SetLineWidth(defaultLineWidth)

	//! ========== DECLARATION BODY ==========
	declarationBodyPadding := 2.0
	declarationTextSize := 9.0
	pdf.SetFont("Arial", "", declarationTextSize)
	declarationBodyX := declarationX + declarationBodyPadding
	declarationBodyY := nextY + 1.5*declarationBodyPadding
	w := pageWidth - 2*defaultPagePadding - 2*declarationBodyPadding
	h := 1.2 * pdf.PointConvert(declarationTextSize)
	declarationText := `    Declaro que não me enquadro no conceito de contribuinte previsto no art. 4º da Lei Complementar nº 87/1996, uma vez que não realizo, com habitualidade ou em volume que caracterize intuito comercial, operações de circulação de mercadoria, ainda que se iniciem no exterior, ou estou dispensado da emissão da nota fiscal por força de legislação tributária vigente, responsabilizando-me, nos termos da lei e a quem de direito, por informações inverídicas.
	  Declaro ainda que não estou postando conteúdo inflamável, explosivo, causador de combustão espontânea, tóxico, corrosivo, gás ou qualquer outro conteúdo que constitua perigo, conforme o art. 13 da Lei Postal nº 6.538/78.`
	pdf.SetXY(declarationBodyX, declarationBodyY)
	pdf.MultiCell(w, h, translator(declarationText), "", "", false)

	//! ========== DECLARATION DATE ==========
	declarationDateX := declarationBodyX + declarationBodyPadding
	declarationDateY := declarationBodyY + 30.0
	declarationDateTextSize := 10.0
	declarationDateString := fmt.Sprintf("%s, %s", city, getCurrentDate())

	pdf.SetFont("Arial", "B", declarationDateTextSize)
	pdf.Text(declarationDateX, declarationDateY, translator(declarationDateString))

	//! ========== DECLARATION SIGNATURE ==========
	declarationSignatureLineX1 := pageWidth/2 + 3*declarationBodyPadding
	declarationSignatureLineX2 := pageWidth - defaultPagePadding - 2*declarationBodyPadding
	declarationSignatureLineY := declarationDateY
	pdf.Line(declarationSignatureLineX1, declarationSignatureLineY, declarationSignatureLineX2, declarationSignatureLineY)

	declarationSignatureString := "Assinatura do Declarante/Remetente"
	declarationSignatureTextX := declarationSignatureLineX1 + 15.5
	declarationSignatureTextY := declarationSignatureLineY + 4.5
	pdf.Text(declarationSignatureTextX, declarationSignatureTextY, translator(declarationSignatureString))

	return y + declarationBlockHeight
}

// ! ==================== OBSERVATIONS BLOCK ====================
func DrawObservations(pdf *gofpdf.Fpdf, x, y float64) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	observationsX := x + defaultPagePadding
	observationsY := y + 2*defaultGridPadding

	pdf.SetLineWidth(defaultBorderLineWidth)
	pdf.Rect(observationsX, observationsY, pageWidth-2*defaultPagePadding, observationBlockHeight, "D")

	observationPadding := 2.5
	observationHeaderTextSize := 12.0

	pdf.SetFont("Arial", "B", observationHeaderTextSize)
	observationHeaderX := observationsX + observationPadding
	observationHeaderY := observationsY + 2*observationPadding

	pdf.Text(observationHeaderX, observationHeaderY, translator("OBSERVAÇÃO:"))

	observationTextX := observationHeaderX - 0.7
	observationTextY := observationHeaderY + 0.5*pdf.PointConvert(observationHeaderTextSize)

	observiationString := "Constitui crime contra a ordem tributária suprimir ou reduzir tributo, ou contribuição social e qualquer acessório (Lei 8.137/90 Art. 1º, V)."
	observationTextSize := 9.0

	observationTextWidth := pageWidth - 2*defaultPagePadding - 2*observationPadding
	observationTextHeight := 1.2 * pdf.PointConvert(observationTextSize)

	pdf.SetFont("Arial", "", observationTextSize)
	fmt.Printf("observationTextX: %.2f\n", observationTextX)
	fmt.Printf("observationTextY: %.2f\n", observationTextY)
	pdf.SetXY(observationTextX, observationTextY)

	pdf.MultiCell(observationTextWidth, observationTextHeight, translator(observiationString), "", "", false)
	return observationsY + observationBlockHeight
}

// ! ==================== HELPER FUNCTIONS ====================
func drawYCenteredText(pdf *gofpdf.Fpdf, x, y, h, fontSize float64, text string) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")
	textY := y + h/2 + pdf.PointConvert(fontSize)/2 - 0.7

	pdf.Text(x, textY, translator(text))
	nextY := y + h
	return nextY
}

func drawXYCenteredText(pdf *gofpdf.Fpdf, x, y, w, h, fontSize float64, text string) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")

	textX := x + w/2 - pdf.GetStringWidth(text)/2
	textY := y + h/2 + pdf.PointConvert(fontSize)/2 - 0.6

	pdf.Text(textX, textY, translator(text))
	nextY := y + h
	return nextY
}

func drawYCenteredXJustifiedRightText(pdf *gofpdf.Fpdf, x, y, w, h, fontSize float64, text string) float64 {
	translator := pdf.UnicodeTranslatorFromDescriptor("")

	textX := x + w - pdf.GetStringWidth(text) - 1.5
	textY := y + h/2 + pdf.PointConvert(fontSize)/2 - 0.5

	pdf.Text(textX, textY, translator(text))
	nextY := y + h
	return nextY
}

func getCurrentDate() string {
	// Set locale to Portuguese (Brazil)
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	// Get the current time in the specified location
	now := time.Now().In(loc)

	// Format the date using the standard Go layout (month will be in English)
	dateStr := now.Format("2 de January de 2006")

	// Replace English month names with Portuguese equivalents
	switch now.Month() {
	case time.January:
		dateStr = strings.Replace(dateStr, "January", "Janeiro", 1)
	case time.February:
		dateStr = strings.Replace(dateStr, "February", "Fevereiro", 1)
	case time.March:
		dateStr = strings.Replace(dateStr, "March", "Março", 1)
	case time.April:
		dateStr = strings.Replace(dateStr, "April", "Abril", 1)
	case time.May:
		dateStr = strings.Replace(dateStr, "May", "Maio", 1)
	case time.June:
		dateStr = strings.Replace(dateStr, "June", "Junho", 1)
	case time.July:
		dateStr = strings.Replace(dateStr, "July", "Julho", 1)
	case time.August:
		dateStr = strings.Replace(dateStr, "August", "Agosto", 1)
	case time.September:
		dateStr = strings.Replace(dateStr, "September", "Setembro", 1)
	case time.October:
		dateStr = strings.Replace(dateStr, "October", "Outubro", 1)
	case time.November:
		dateStr = strings.Replace(dateStr, "November", "Novembro", 1)
	case time.December:
		dateStr = strings.Replace(dateStr, "December", "Dezembro", 1)
	}

	return dateStr
}

func formatCEP(cep string) string {
	if len(cep) != 8 {
		// Handle invalid CEP length appropriately
		return "Invalid CEP"
	}

	return cep[:5] + "-" + cep[5:]
}

func formatCPFCNPJ(cpfCnpj int) string {
	// Convert the number to a string
	str := fmt.Sprintf("%d", cpfCnpj)

	// Determine if it's CNPJ or CPF based on length and pad with leading zeros
	switch {
	case len(str) > 11: // CNPJ, 13 or 14 digits
		str = fmt.Sprintf("%014d", cpfCnpj)
		return fmt.Sprintf("%s.%s.%s/%s-%s", str[0:2], str[2:5], str[5:8], str[8:12], str[12:14])
	case len(str) == 11, len(str) == 10: // CPF, 11 digits
		str = fmt.Sprintf("%011d", cpfCnpj)
		return fmt.Sprintf("%s.%s.%s-%s", str[0:3], str[3:6], str[6:9], str[9:11])
	default:
		// Return the original string if it's neither CNPJ nor CPF
		return fmt.Sprintf("%d", cpfCnpj)
	}
}

func formatCPFCNPJV2(cpfCnpj string) string {
	// Convert the number to a string
	str := fmt.Sprintf("%s", cpfCnpj)

	// Determine if it's CNPJ or CPF based on length and pad with leading zeros
	switch {
	case len(str) > 11: // CNPJ, 13 or 14 digits
		str = fmt.Sprintf("%014s", cpfCnpj)
		return fmt.Sprintf("%s.%s.%s/%s-%s", str[0:2], str[2:5], str[5:8], str[8:12], str[12:14])
	case len(str) == 11, len(str) == 10: // CPF, 11 digits
		str = fmt.Sprintf("%011s", cpfCnpj)
		return fmt.Sprintf("%s.%s.%s-%s", str[0:3], str[3:6], str[6:9], str[9:11])
	default:
		// Return the original string if it's neither CNPJ nor CPF
		return fmt.Sprintf("%s", cpfCnpj)
	}
}
