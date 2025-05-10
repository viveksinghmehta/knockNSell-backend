package main

import (
	"context"
	"knockNSell/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	router := gin.Default()
	router.POST("/sendotp", routes.Sendotp)
	ginLambda = ginadapter.NewV2(router)
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
}
