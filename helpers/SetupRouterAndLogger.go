package helper

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetUpRouterAndLogger(environment string) *gin.Engine {
	// Set up Logrus with JSON formatter
	log.SetFormatter(&log.JSONFormatter{})

	// Set up Logflare hook
	logflareHook := NewLogflareHook(os.Getenv("LOGFLARE_API_KEY"), os.Getenv("LOGFLARE_SOURCE_ID"))

	// Set up Slack hook
	slackHook := NewSlackHook(os.Getenv("SLACK_WEBHOOK_URL"))
	log.AddHook(slackHook)

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
	router.Use(gin.Recovery())

	return router
}
