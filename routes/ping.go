package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PingServer(c *gin.Context) {
	log.Infof("Request sent to the server :Ping")
	c.Status(http.StatusOK)
}
