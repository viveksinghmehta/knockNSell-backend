package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// SlackHook is a Logrus hook to send error logs to Slack
type SlackHook struct {
	WebhookURL string
}

// NewSlackHook creates a new Slack hook
func NewSlackHook(webhookURL string) *SlackHook {
	return &SlackHook{
		WebhookURL: webhookURL,
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
			"value": entry.Level.String(),
			"short": "true",
		},
		{
			"title": "Timestamp",
			"value": entry.Time.Format(time.RFC3339),
			"short": "true",
		},
	}

	// Add custom fields from entry.Data
	for key, value := range entry.Data {
		fields = append(fields, map[string]string{
			"title": key,
			"value": fmt.Sprintf("%v", value),
			"short": "true",
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
