// Package config loads environment-backed configuration for the agent app.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/joho/godotenv"
)

// Config is the complete runtime configuration for the agent app.
type Config struct {
	Provider llm.ProviderType
	Agent    AgentConfig
	OpenAI   OpenAIConfig
}

// AgentConfig contains agent loop limits.
type AgentConfig struct {
	MaxRounds   int
	ToolTimeout time.Duration
}

// OpenAIConfig contains settings for the OpenAI LLM provider.
type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// LoadEnv loads local environment variables from a .env file when present.
func LoadEnv() error {
	return godotenv.Load()
}

// Load reads environment variables and validates the selected provider config.
func Load() (Config, error) {
	cfg := Config{
		Provider: envProviderOrDefault("LLM_PROVIDER", llm.ProviderTypeOpenAI),
		Agent: AgentConfig{
			MaxRounds:   envIntOrDefault("AGENT_MAX_ROUNDS", 0),
			ToolTimeout: envDurationOrDefault("AGENT_TOOL_TIMEOUT", 15*time.Second),
		},
		OpenAI: OpenAIConfig{
			APIKey:  strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
			Model:   envOrDefault("OPENAI_MODEL", "gpt-5.4-mini"),
			BaseURL: strings.TrimRight(envOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"), "/"),
		},
	}

	if cfg.Provider == llm.ProviderTypeOpenAI && cfg.OpenAI.APIKey == "" {
		return Config{}, errors.New("OPENAI_API_KEY is required when LLM_PROVIDER=openai")
	}

	return cfg, nil
}

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
