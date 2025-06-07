package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	logger "knockNSell/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests and errors
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read body (if possible)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restore io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		headersJSON, _ := json.Marshal(c.Request.Header)

		msg := logger.GetLogMessage(c.Request.Context())
		extraFields := logger.GetExtraFields(c.Request.Context())

		if msg == "" {
			if statusCode >= 400 {
				msg = "HTTP Error"
			} else {
				msg = "HTTP Request"
			}
		}

		args := []any{
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("headers", string(headersJSON)),
			slog.String("body", string(bodyBytes)),
			slog.Float64("duration_seconds", duration.Seconds()),
			slog.String("error", errMsg),
		}
		// add extra fields as separate fields:
		for k, v := range extraFields {
			args = append(args, slog.Any(k, v))
		}

		if statusCode >= 400 {
			slog.Error(msg, args...)

		} else {
			slog.Info(msg, args...)
		}
	}
}
