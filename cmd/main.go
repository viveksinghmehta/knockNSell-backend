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
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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

func setUpRouterAndLogger(environment string) *gin.Engine {
	// Set up Logrus with JSON formatter
	log.SetFormatter(&log.JSONFormatter{})

	// Set up Logflare hook (replace with your Source ID and API Key)
	logflareHook := helper.NewLogflareHook(os.Getenv("LOGFLARE_API_KEY"), os.Getenv("lOGFLARE_SOURCE_ID"))
	log.AddHook(logflareHook)
	// 3. Wire up Gin with your handlers
	router := gin.New()

	// Custom middleware to log requests
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.WithFields(log.Fields{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration":    duration.String(),
			"environment": environment,
			"ip":          c.ClientIP(),
		}).Info("Request completed")
	})

	return router
}

func init() {
	// Only for development remove it in PROD
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Use this to get the mode :- release/debug
	mode := os.Getenv("GIN_MODE")

	initDB()

	if mode == "release" {

		router := setUpRouterAndLogger("PROD")
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		router.POST("/login", server.LoginUser)
		router.POST("/signup", server.SignUpUser)
		ginLambda = ginadapter.NewV2(router)
	}
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "debug" {
		router := setUpRouterAndLogger("DEV")
		server := routes.NewServer(queries)
		router.GET("/ping", routes.PingServer)
		router.POST("/sendotp", server.Sendotp)
		router.POST("/verifyotp", server.VerifyOTP)
		router.POST("/login", server.LoginUser)
		router.POST("/signup", server.SignUpUser)
		router.Run()
	} else {
		lambda.Start(handleRequest) // start Lambda :contentReference[oaicite:9]{index=9}
	}
}
