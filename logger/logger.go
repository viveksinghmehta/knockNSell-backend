package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type contextKey string

const (
	logMessageKey  = contextKey("log_message")
	extraFieldsKey = contextKey("extra_fields")
)

// SetLogMessageAndFields sets both the log message and additional fields in the request context
func SetLogMessageAndFields(ctx context.Context, message string, fields map[string]any) context.Context {
	ctx = SetLogMessage(ctx, message)
	if fields != nil {
		ctx = SetExtraFields(ctx, fields)
	}
	return ctx
}

// SetLogMessage sets a custom log message in the request context
func SetLogMessage(ctx context.Context, message string) context.Context {
	return context.WithValue(ctx, logMessageKey, message)
}

// GetLogMessage retrieves the custom log message from the request context
func GetLogMessage(ctx context.Context) string {
	if msg, ok := ctx.Value(logMessageKey).(string); ok {
		return msg
	}
	return ""
}

// SetExtraFields sets additional fields in the request context for logging
func SetExtraFields(ctx context.Context, fields map[string]any) context.Context {
	return context.WithValue(ctx, extraFieldsKey, fields)
}

// GetExtraFields retrieves additional fields from the request context for logging
func GetExtraFields(ctx context.Context) map[string]any {
	if fields, ok := ctx.Value(extraFieldsKey).(map[string]any); ok {
		return fields
	}
	return nil
}

type DiscordPayload struct {
	Content string `json:"content"`
}

func sendToDiscordStructured(obj any) {
	// Marshal your object
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return
	}

	payload := DiscordPayload{Content: "```json\n" + string(formatted) + "\n```"}
	data, _ := json.Marshal(payload)

	http.Post(os.Getenv("DISCORD_WEBHOOK_URL"), "application/json", bytes.NewBuffer(data))
}

func sendToBetterStack(message string) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", os.Getenv("BETTERSTACK_URL"), bytes.NewBuffer([]byte(message)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("BETTERSTACK_TOKEN"))

	client.Do(req)
}

type CustomHandler struct {
	h   slog.Handler
	env string
}

func (ch *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ch.h.Enabled(ctx, level)
}

func (ch *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	ch.h.Handle(ctx, record)
	err := ch.h.Handle(ctx, record)

	if record.Level >= slog.LevelError {
		if ch.env != "debug" {
			// Collect attrs
			attrMap := make(map[string]any)
			record.Attrs(func(a slog.Attr) bool {
				attrMap[a.Key] = a.Value.Any()
				return true
			})

			// Build full structured payload
			errorPayload := map[string]any{
				"time":    record.Time.Format(time.RFC3339),
				"level":   record.Level.String(),
				"message": record.Message,
				"attrs":   attrMap,
			}

			// Send to Discord
			go sendToDiscordStructured(errorPayload)

			// Send to BetterStack as plain text
			plainMessage, _ := json.MarshalIndent(errorPayload, "", "  ")
			go sendToBetterStack(string(plainMessage))
		}
	}

	return err
}

func (ch *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{h: ch.h.WithAttrs(attrs), env: ch.env}
}

func (ch *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{h: ch.h.WithGroup(name), env: ch.env}
}

func SetupLogger() {
	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = "debug"
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(&CustomHandler{h: h, env: env})
	slog.SetDefault(logger)
}
