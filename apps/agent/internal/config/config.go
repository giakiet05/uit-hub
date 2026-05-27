// Package config loads environment-backed configuration for the agent app.
package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/joho/godotenv"
)

// Config is the complete runtime configuration for the agent app.
type Config struct {
	Provider llm.ProviderType
	Agent    AgentConfig
	Memory   MemoryConfig
	MCP      MCPConfig
	OpenAI   OpenAIConfig
	Tool     ToolConfig
}

// AgentConfig contains agent loop limits.
type AgentConfig struct {
	Type        agent.Type
	MaxRounds   int
	ToolTimeout time.Duration
}

// MemoryConfig contains long-term memory settings.
type MemoryConfig struct {
	Path string
}

// ToolConfig contains tool registry settings.
type ToolConfig struct {
	RuntimeLimit int
}

// OpenAIConfig contains settings for the OpenAI LLM provider.
type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// LoadEnv loads local environment variables from the nearest known agent .env file when present.
func LoadEnv() error {
	candidates := []string{
		".env",
		filepath.Join("apps", "agent", ".env"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return godotenv.Load(candidate)
		}
	}
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
			Type:        agent.Type(envOrDefault("AGENT_TYPE", agent.TypeReAct.String())),
			MaxRounds:   envIntOrDefault("AGENT_MAX_ROUNDS", 0),
			ToolTimeout: envDurationOrDefault("AGENT_TOOL_TIMEOUT", 15*time.Second),
		},
		Memory: MemoryConfig{
			Path: envOrDefault("MEMORY_PATH", "memory"),
		},
		MCP: mcpConfig,
		Tool: ToolConfig{
			RuntimeLimit: envIntOrDefault("TOOL_RUNTIME_LIMIT", 20),
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
