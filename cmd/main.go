package main

import (
	"context"
	"database/sql"
	db "knockNSell/db/gen"
	"knockNSell/routes"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambdaV2
var queries *db.Queries

func init() {

	// Only for development remove it in PROD
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 1. Open the database once (Lambda warm starts reuse this)
	dsn := os.Getenv("DATABASE_URL")
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	// 2. Instantiate sqlc queries
	queries = db.New(sqldb)

	// 3. Wire up Gin with your handlers
	// Commented the code to run in development
	// router := gin.Default()
	// server := routes.NewServer(queries)
	// router.GET("/ping", routes.PingServer)
	// router.POST("/sendotp", server.Sendotp)
	// ginLambda = ginadapter.NewV2(router)
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Commented this code for Development
	// lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
	router := gin.Default()
	server := routes.NewServer(queries)
	router.GET("/ping", routes.PingServer)
	router.POST("/sendotp", server.Sendotp)
	router.Run()
}
