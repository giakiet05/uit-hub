package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SystemPrompt string
	Provider     string
	OpenAI       OpenAIConfig
}

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func LoadEnv() error {
	return godotenv.Load()
}

func Load() (Config, error) {
	cfg := Config{
		SystemPrompt: "You are a helpful agent.",
		Provider:     envOrDefault("LLM_PROVIDER", "openai"),
		OpenAI: OpenAIConfig{
			APIKey:  strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
			Model:   envOrDefault("OPENAI_MODEL", "gpt-4.1-mini"),
			BaseURL: strings.TrimRight(envOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"), "/"),
		},
	}

	if cfg.Provider == "openai" && cfg.OpenAI.APIKey == "" {
		return Config{}, errors.New("OPENAI_API_KEY is required when LLM_PROVIDER=openai")
	}

	return cfg, nil
}

func DefaultForTest() Config {
	return Config{
		SystemPrompt: "You are a helpful agent.",
		Provider:     "echo",
	}
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
