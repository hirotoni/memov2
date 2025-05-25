package common

import (
	"log/slog"
	"os"
)

// Config holds logger configuration
type LoggerConfig struct {
	Level  string
	Format string
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(cfg LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	default:
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}

// DefaultLogger creates a logger with default configuration
func DefaultLogger() *slog.Logger {
	return NewLogger(LoggerConfig{
		Level:  "info",
		Format: "text",
	})
}

