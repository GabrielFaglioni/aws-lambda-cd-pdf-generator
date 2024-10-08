package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FavorDespaches/content-declaration-pdf-generator/pkg/helpers"
	"github.com/FavorDespaches/content-declaration-pdf-generator/pkg/types"
	"github.com/aws/aws-lambda-go/events"
)

func SolicitarDeclaracaoConteudoLocalV2(solicitarDeclaracaoConteudo types.SolicitarDeclaracaoConteudo) error {
	fmt.Println("\n\n========== INICIANDO LAMBDA ==========")

	return helpers.GenerateContentDeclarationPDFLocal(solicitarDeclaracaoConteudo)
}

func SolicitarDeclaracaoConteudo(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("\n\n========== INICIANDO LAMBDA ==========")
	var solicitarDeclaracaoConteudo types.SolicitarDeclaracaoConteudo
	err := json.Unmarshal([]byte(req.Body), &solicitarDeclaracaoConteudo)

	if err != nil {
		errText := fmt.Sprintf("Erro no Parser do JSON: %s", err.Error())

		errorBody := ErrorResponse{
			Message: errText,
		}
		return ApiResponse(http.StatusBadRequest, errorBody)
	}

	base64String, err := helpers.GenerateContentDeclaration(solicitarDeclaracaoConteudo)

	if err != nil {
		errText := fmt.Sprintf("Erro GenerateLabelsPDF: %s", err.Error())

		errorBody := ErrorResponse{
			Message: errText,
		}
		ApiResponse(http.StatusBadRequest, errorBody)
	}

	successBody := SuccessResponse{
		StringBase64: base64String,
	}

	return ApiResponse(http.StatusOK, successBody)
}

func HandleUnsupportedMethod(method string) (*events.APIGatewayProxyResponse, error) {
	errorBody := ErrorResponse{
		Message: fmt.Sprintf("Método inválido: %s. Utilize somente POST", method),
	}

	return ApiResponse(http.StatusBadRequest, errorBody)
}
