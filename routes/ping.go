package routes

import (
	"net/http"
	"time"

	helper "knockNSell/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PingServer(c *gin.Context) {
	log.Info("Request sent to the server :Ping")
	c.Status(http.StatusOK)
}

func SendError(c *gin.Context) {
	start := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"status":  400,
		"message": "successfully send an Error.",
	})
	log.WithFields(helper.GetExtraFieldsForSlackLog(c, start)).Error("Check the error 🚨")
}
