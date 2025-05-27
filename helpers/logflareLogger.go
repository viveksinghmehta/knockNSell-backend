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

	// Determine color based on log level
	var color string
	switch entry.Level {
	case log.InfoLevel:
		color = "#00FF00" // Green
	case log.ErrorLevel:
		color = "#FF0000" // Red
	case log.WarnLevel:
		color = "#FFA500" // Orange
	case log.DebugLevel:
		color = "#0000FF" // Blue
	default:
		color = "#FFFFFF" // White for other levels (e.g., trace)
	}
	// Prepare the log payload
	payload := map[string]interface{}{
		"message": entry.Message,
		"metadata": map[string]interface{}{
			"Level":     entry.Level.String(),
			"Timestamp": entry.Time.Format(time.RFC850),
			"color":     color,
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
