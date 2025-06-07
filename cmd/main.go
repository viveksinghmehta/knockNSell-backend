package main

import (
	"context"
	"database/sql"
	db "knockNSell/db/gen"
	helper "knockNSell/helpers"
	"knockNSell/routes"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	log "github.com/sirupsen/logrus"
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

	// Use this to get the mode :- release/debug
	mode := os.Getenv("GIN_MODE")

	initDB()

	if mode == "release" {
		router := helper.SetUpRouterAndLogger("PROD")
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		router.POST("/login", server.LoginUser)
		router.POST("/updateProfile", server.UpdateProfile)
		router.POST("/signup", server.SignUpUser)
		router.POST("/error", routes.SendError)
		ginLambda = ginadapter.NewV2(router)
	}
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "debug" {
		router := helper.SetUpRouterAndLogger("PROD")
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		router.POST("/login", server.LoginUser)
		router.POST("/updateProfile", server.UpdateProfile)
		router.POST("/signup", server.SignUpUser)
		router.POST("/error", routes.SendError)
		router.Run()
	} else {
		lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
	}
}
