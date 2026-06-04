// Package config loads environment-backed configuration for the agent app.
package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/joho/godotenv"
)

// Config is the complete runtime configuration for the agent app.
type Config struct {
	Registry   *llm.Registry
	Agent      AgentConfig
	Memory     MemoryConfig
	MCP        MCPConfig
	Tool       ToolConfig
	Compaction CompactionConfig
}

// AgentConfig contains agent loop limits.
type AgentConfig struct {
	Type            agent.Type
	MaxRounds       int
	ToolTimeout     time.Duration
	ConcurrentTools bool
}

// MemoryConfig contains long-term memory settings.
type MemoryConfig struct {
	Path string
}

// ToolConfig contains tool registry settings.
type ToolConfig struct {
	RuntimeLimit int
}

// CompactionConfig contains lightweight conversation compaction settings.
type CompactionConfig struct {
	ToolResultBudgetEnabled           bool
	ToolResultBudgetMaxChars          int
	SnipEnabled                       bool
	SnipTriggerRatio                  float64
	MicrocompactEnabled               bool
	MicrocompactTriggerRatio          float64
	MicrocompactMinChars              int
	MicrocompactKeepRecentToolResults int
	ContextCollapseEnabled            bool
	ContextCollapseTriggerRatio       float64
	ContextCollapseTargetRatio        float64
	ContextCollapseKeepRecentTurns    int
	ContextCollapseMaxSegments        int
	ContextCollapseMaxConcurrency     int
	ContextCollapseMinSegmentTokens   int
	AutoCompactEnabled                bool
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

	factoryPath := os.Getenv("MODELS_CONFIG_PATH")
	if factoryPath == "" {
		if _, err := os.Stat("models.yaml"); err == nil {
			factoryPath = "models.yaml"
		} else {
			factoryPath = filepath.Join("apps", "agent", "models.yaml")
		}
	}

	logger := slog.Default() // You could also inject this or set up a dummy one for loading.
	registry, err := LoadLLMRegistry(factoryPath, logger)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Registry: registry,
		Agent: AgentConfig{
			Type:            agent.Type(envOrDefault("AGENT_TYPE", string(agent.TypeReAct))),
			MaxRounds:       envIntOrDefault("AGENT_MAX_ROUNDS", 0),
			ToolTimeout:     envDurationOrDefault("AGENT_TOOL_TIMEOUT", 15*time.Minute),
			ConcurrentTools: envBoolOrDefault("AGENT_CONCURRENT_TOOLS", true),
		},
		Memory: MemoryConfig{
			Path: envOrDefault("MEMORY_PATH", "memory"),
		},
		MCP: mcpConfig,
		Tool: ToolConfig{
			RuntimeLimit: envIntOrDefault("TOOL_RUNTIME_LIMIT", 20),
		},
		Compaction: CompactionConfig{
			ToolResultBudgetEnabled:           envBoolOrDefault("COMPACTION_TOOL_RESULT_BUDGET_ENABLED", true),
			ToolResultBudgetMaxChars:          envIntOrDefault("COMPACTION_TOOL_RESULT_BUDGET_MAX_CHARS", 12000),
			SnipEnabled:                       envBoolOrDefault("COMPACTION_SNIP_ENABLED", true),
			SnipTriggerRatio:                  envFloatOrDefault("COMPACTION_SNIP_TRIGGER_RATIO", 0.95),
			MicrocompactEnabled:               envBoolOrDefault("COMPACTION_MICROCOMPACT_ENABLED", true),
			MicrocompactTriggerRatio:          envFloatOrDefault("COMPACTION_MICROCOMPACT_TRIGGER_RATIO", 0.70),
			MicrocompactMinChars:              envIntOrDefault("COMPACTION_MICROCOMPACT_MIN_CHARS", 12000),
			MicrocompactKeepRecentToolResults: envIntOrDefault("COMPACTION_MICROCOMPACT_KEEP_RECENT_TOOL_RESULTS", 8),
			ContextCollapseEnabled:            envBoolOrDefault("COMPACTION_CONTEXT_COLLAPSE_ENABLED", false),
			ContextCollapseTriggerRatio:       envFloatOrDefault("COMPACTION_CONTEXT_COLLAPSE_TRIGGER_RATIO", 0.80),
			ContextCollapseTargetRatio:        envFloatOrDefault("COMPACTION_CONTEXT_COLLAPSE_TARGET_RATIO", 0.65),
			ContextCollapseKeepRecentTurns:    envIntOrDefault("COMPACTION_CONTEXT_COLLAPSE_KEEP_RECENT_TURNS", 3),
			ContextCollapseMaxSegments:        envIntOrDefault("COMPACTION_CONTEXT_COLLAPSE_MAX_SEGMENTS", 3),
			ContextCollapseMaxConcurrency:     envIntOrDefault("COMPACTION_CONTEXT_COLLAPSE_MAX_CONCURRENCY", 2),
			ContextCollapseMinSegmentTokens:   envIntOrDefault("COMPACTION_CONTEXT_COLLAPSE_MIN_SEGMENT_TOKENS", 1500),
			AutoCompactEnabled:                envBoolOrDefault("COMPACTION_AUTO_COMPACT_ENABLED", false),
		},
	}

	if cfg.Agent.Type == "" {
		return Config{}, errors.New("AGENT_TYPE is required")
	}

	return cfg, nil
}
