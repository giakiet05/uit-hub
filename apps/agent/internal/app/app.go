// Package app wires configuration, providers, tools, sessions, and the TUI.
package app

import (
	"context"
	"io"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
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
		"llm_provider", cfg.Provider,
		"openai_model", cfg.OpenAI.Model,
		"openai_base_url", cfg.OpenAI.BaseURL,
		"agent_max_rounds", cfg.Agent.MaxRounds,
		"agent_tool_timeout", cfg.Agent.ToolTimeout.String(),
	)
	logger.DebugContext(ctx, "Starting agent app", "llm_provider", cfg.Provider)

	provider, err := newProvider(cfg, logger)
	if err != nil {
		return err
	}

	session := runtime.NewSession()
	prompts := prompt.NewBuilder(prompt.DefaultSystemPrompt())
	tools, err := newToolRegistry()
	if err != nil {
		logger.DebugContext(ctx, "Tool registry initialization failed", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Tool registry initialized", "tool_count", len(tools.Definitions()))

	runtimeAgent := agent.NewReActAgent(agent.ReActConfig{
		Provider:    provider,
		Prompts:     prompts,
		Tools:       tools,
		Logger:      logger,
		MaxRounds:   cfg.Agent.MaxRounds,
		ToolTimeout: cfg.Agent.ToolTimeout,
	})

	runner := tui.NewRunner(session, runtimeAgent, stdin, stdout, logBuffer, initialPrompt)
	if err := runner.Run(ctx); err != nil {
		logger.DebugContext(ctx, "Agent TUI stopped with error", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent TUI stopped")
	return nil
}
