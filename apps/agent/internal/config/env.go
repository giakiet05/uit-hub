package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

// envProviderOrDefault reads a ProviderType environment variable.
func envProviderOrDefault(key string, fallback llm.ProviderType) llm.ProviderType {
	return llm.ProviderType(envOrDefault(key, string(fallback)))
}

// envOrDefault reads a trimmed string environment variable.
func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

// envIntOrDefault reads an integer environment variable.
func envIntOrDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// envDurationOrDefault reads a time.Duration environment variable.
func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
