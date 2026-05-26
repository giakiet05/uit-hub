// Package logging configures slog loggers for the agent runtime.
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a text slog logger whose level is controlled by LOG_LEVEL.
func NewLogger(output io.Writer) *slog.Logger {
	logLevel := slog.LevelInfo

	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(output, opts)
	logger := slog.New(handler)

	if logLevel == slog.LevelDebug {
		logger.Info("Debug logging enabled")
	} else {
		logger.Info("Info logging enabled")
	}

	return logger
}

// NewNopLogger creates a logger that discards all records.
func NewNopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
