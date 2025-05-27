package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// SlackHook is a Logrus hook to send error logs to Slack
type SlackHook struct {
	WebhookURL  string
	AppName     string // Optional: for custom app name field
	Environment string // Optional: for custom environment field
}

// NewSlackHook creates a new Slack hook
func NewSlackHook(webhookURL, appName, environment string) *SlackHook {
	return &SlackHook{
		WebhookURL:  webhookURL,
		AppName:     appName,
		Environment: environment,
	}
}

// Levels returns the log levels that trigger this hook
func (hook *SlackHook) Levels() []log.Level {
	return []log.Level{log.ErrorLevel} // Only send error logs to Slack
}

// Fire sends the error log entry to Slack
func (hook *SlackHook) Fire(entry *log.Entry) error {
	// Initialize the fields slice
	fields := []map[string]string{
		{
			"title": "Level",
			"value": strings.ToUpper(entry.Level.String()),
			"short": "true",
		},
		{
			"title": "Timestamp",
			"value": entry.Time.Format(time.RFC850),
			"short": "true",
		},
		{
			"title": "App",
			"value": hook.AppName,
			"short": "true",
		},
		{
			"title": "Environment",
			"value": hook.Environment,
			"short": "true",
		},
	}

	// Add all fields from entry.Data
	for key, value := range entry.Data {
		fields = append(fields, map[string]string{
			"title": key,
			"value": fmt.Sprintf("%v", value),
			"short": "false",
		})
	}

	// Prepare the Slack payload
	payload := map[string]interface{}{
		"text": entry.Message,
		"attachments": []map[string]interface{}{
			{
				"color":  "#FF0000", // Red for errors
				"fields": fields,
			},
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", hook.WebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack returned status: %d", resp.StatusCode)
	}

	return nil
}

func GetExtraFieldsForSlackLog(c *gin.Context, startTime time.Time) log.Fields {
	duration := time.Since(startTime)

	// Read request body (copying for safety)
	var requestBody string
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		requestBody = string(bodyBytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Capture request headers
	headerBytes, _ := json.MarshalIndent(c.Request.Header, "", "  ")

	// Attempt to get response body from context (assuming it was set previously)
	responseBody := c.Writer.Header().Get("X-Response-Body") // Custom: must be set by middleware if needed

	return log.Fields{
		"Path":           c.Request.URL.Path,
		"Status code":    c.Writer.Status(),
		"IP address":     c.Request.RemoteAddr,
		"Method":         c.Request.Method,
		"Duration":       duration.String(),
		"Host":           c.Request.Host,
		"Request Header": string(headerBytes),
		"Request Body":   requestBody,
		"Response Body":  responseBody,
	}
}
