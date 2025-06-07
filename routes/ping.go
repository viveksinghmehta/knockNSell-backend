package routes

import (
	logger "knockNSell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingServer(c *gin.Context) {
	c.Status(http.StatusOK)
}

func SendError(c *gin.Context) {
	c.Request = c.Request.WithContext(
		logger.SetLogMessage(c.Request.Context(), "✅ Successfully send error log"),
	)
	c.JSON(400, gin.H{
		"status":  400,
		"message": "successfully send an Error.",
	})
}
