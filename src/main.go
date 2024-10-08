package main

import (
	"strings"

	"github.com/FavorDespaches/content-declaration-pdf-generator/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	method := strings.ToUpper(req.HTTPMethod)

	if method != "POST" {
		return handlers.HandleUnsupportedMethod(method)
	}

	return handlers.SolicitarDeclaracaoConteudo(req)
}
