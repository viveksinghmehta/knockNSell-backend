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
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambdaV2
var queries *db.Queries

func init() {
	// Only for development remove it in PROD
	mode := os.Getenv("GIN_MODE")

	if mode == "debug" {
		// Only for development remove it in PROD
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
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

	if mode == "release" {
		// 3. Wire up Gin with your handlers
		router := gin.Default()
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		ginLambda = ginadapter.NewV2(router)
	}
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "debug" {
		router := gin.Default()
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		router.Run()
	} else {
		lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
	}
}
