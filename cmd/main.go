package main

import (
	"context"
	"database/sql"
	"fmt"
	db "knockNSell/db/gen"
	logger "knockNSell/logger"
	middleWare "knockNSell/middleware"
	"knockNSell/routes"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambdaV2
var queries *db.Queries

func initDB() {
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
}

func init() {
	// Remove this code for Production
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Use this to get the mode :- release/debug
	mode := os.Getenv("GIN_MODE")

	initDB()

	logger.SetupLogger()

	if mode == "release" {

		router := routes.SetupRouter()
		router.Use(middleWare.LoggingMiddleware())
		routes.RegisterRoutes(router, queries)
		ginLambda = ginadapter.NewV2(router)
	}
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "debug" {
		router := routes.SetupRouter()
		router.Use(middleWare.LoggingMiddleware())
		routes.RegisterRoutes(router, queries)
		router.Run()
	} else {
		lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
	}
}
