package main

import (
	"knockNSell/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/sendotp", routes.Sendotp)
	router.Run()
}
