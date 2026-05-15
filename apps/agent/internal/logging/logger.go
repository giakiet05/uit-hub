package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

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

func NewNopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
