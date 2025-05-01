package main

import (
	"knockNSell/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	router := gin.Default()
	router.POST("/sendotp", routes.Sendotp)
	router.Run()
}
