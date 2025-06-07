package routes

import (
	db "knockNSell/db/gen"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

func RegisterRoutes(router *gin.Engine, queries *db.Queries) {
	server := NewServer(queries)

	router.GET("/ping", PingServer)
	router.POST("/sendotp", server.Sendotp)
	router.POST("/verifyotp", server.VerifyOTP)
	router.POST("/login", server.LoginUser)
	router.POST("/updateProfile", server.UpdateProfile)
	router.POST("/signup", server.SignUpUser)
	router.POST("/error", SendError)
}
