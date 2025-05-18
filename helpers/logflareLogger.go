package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogflareHook is a Logrus hook to send logs to Logflare
type LogflareHook struct {
	APIKey   string
	SourceID string
}

// NewLogflareHook creates a new Logflare hook
func NewLogflareHook(apiKey, sourceID string) *LogflareHook {
	return &LogflareHook{
		APIKey:   apiKey,
		SourceID: sourceID,
	}
}

// Levels returns the log levels that trigger this hook
func (hook *LogflareHook) Levels() []log.Level {
	return log.AllLevels // Send all log levels to Logflare
}

// Fire sends the log entry to Logflare
func (hook *LogflareHook) Fire(entry *log.Entry) error {
	// Prepare the log payload
	payload := map[string]interface{}{
		"message": entry.Message,
		"metadata": map[string]interface{}{
			"level":     entry.Level.String(),
			"timestamp": entry.Time.Format(time.RFC3339),
			"fields":    entry.Data,
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.logflare.app/logs/json?source="+hook.SourceID, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", hook.APIKey)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Logflare returned status: %d", resp.StatusCode)
	}

	return nil
}
