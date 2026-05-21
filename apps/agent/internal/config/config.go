// Package config loads environment-backed configuration for the agent app.
package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/joho/godotenv"
)

// Config is the complete runtime configuration for the agent app.
type Config struct {
	Provider llm.ProviderType
	Agent    AgentConfig
	MCP      MCPConfig
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
	mcpConfig, err := loadMCPConfig(envOrDefault("MCP_CONFIG_PATH", "mcp.yaml"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Provider: envProviderOrDefault("LLM_PROVIDER", llm.ProviderTypeOpenAI),
		Agent: AgentConfig{
			MaxRounds:   envIntOrDefault("AGENT_MAX_ROUNDS", 0),
			ToolTimeout: envDurationOrDefault("AGENT_TOOL_TIMEOUT", 15*time.Second),
		},
		MCP: mcpConfig,
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
