package main

import (
	"context"
	"knockNSell/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	router := gin.Default()
	router.POST("/sendotp", routes.Sendotp)
	router.Run()
	ginLambda = ginadapter.New(router)
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy the request through Gin
	return ginLambda.ProxyWithContext(ctx, req) // proxy & return response :contentReference[oaicite:8]{index=8}
}

func main() {
	lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
}
