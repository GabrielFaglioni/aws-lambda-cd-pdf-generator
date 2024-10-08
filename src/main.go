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
	switch strings.ToUpper(req.HTTPMethod) {
	case "GET":
		return handlers.HandleUnsupportedMethod()
	case "POST":
		return handlers.SolicitarDeclaracaoConteudoV2(req) //Para dar update para a V2, basta substituir para 'SolicitarDeclaracaoConteudoV2'
	case "PUT":
		return handlers.HandleUnsupportedMethod()
	case "DELETE":
		return handlers.HandleUnsupportedMethod()
	default:
		return handlers.HandleUnsupportedMethod()
	}
}
