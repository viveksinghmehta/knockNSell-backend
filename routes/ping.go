package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PingServer(c *gin.Context) {
	log.Info("Request sent to the server :Ping")
	c.Status(http.StatusOK)
}

func SendError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  400,
		"message": "successfully send an Error.",
	})
	log.Error("Saved the error log.")
}
