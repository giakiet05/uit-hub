// Package app wires configuration, providers, tools, sessions, and the TUI.
package app

import (
	"context"
	"io"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tui"
)

// Run starts the agent application in TUI mode.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	logBuffer := tui.NewLogBuffer(500)
	initialPrompt := strings.TrimSpace(strings.Join(args, " "))

	err := config.LoadEnv()
	logger := logging.NewLogger(logBuffer)
	if err != nil {
		logger.DebugContext(ctx, "No .env file loaded", "error", err)
	} else {
		logger.DebugContext(ctx, ".env file loaded")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.DebugContext(ctx, "Agent config load failed", "error", err)
		return err
	}
	logger.DebugContext(
		ctx,
		"Agent config loaded",
		"agent_type", cfg.Agent.Type,
		"llm_provider", cfg.Provider,
		"openai_model", cfg.OpenAI.Model,
		"openai_base_url", cfg.OpenAI.BaseURL,
		"agent_max_rounds", cfg.Agent.MaxRounds,
		"agent_tool_timeout", cfg.Agent.ToolTimeout.String(),
		"memory_path", cfg.Memory.Path,
		"tool_runtime_limit", cfg.Tool.RuntimeLimit,
		"mcp_config_path", cfg.MCP.ConfigPath,
		"mcp_server_count", len(cfg.MCP.Servers),
	)
	logger.DebugContext(ctx, "Starting agent app", "llm_provider", cfg.Provider)

	runtime, err := NewRuntime(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer runtime.Close(ctx)

	runner := tui.NewRunner(runtime.Session, stdin, stdout, logBuffer, initialPrompt)
	if err := runner.Run(ctx); err != nil {
		logger.DebugContext(ctx, "Agent TUI stopped with error", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent TUI stopped")
	return nil
}
