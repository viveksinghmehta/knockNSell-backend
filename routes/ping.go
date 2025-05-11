package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingServer(c *gin.Context) {
	c.Status(http.StatusOK)
}
