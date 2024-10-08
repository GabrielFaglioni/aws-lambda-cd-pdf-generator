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

	switch method {
	case "GET":
		return handlers.HandleUnsupportedMethod(method)
	case "POST":
		return handlers.SolicitarDeclaracaoConteudoV2(req)
	case "PUT":
		return handlers.HandleUnsupportedMethod(method)
	case "DELETE":
		return handlers.HandleUnsupportedMethod(method)
	default:
		return handlers.HandleUnsupportedMethod(method)
	}
}
