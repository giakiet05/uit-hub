// Package app wires configuration, providers, tools, sessions, and the TUI.
package app

import (
	"context"
	"io"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
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
		"memory_path", cfg.Memory.Path,
		"tool_runtime_limit", cfg.Tool.RuntimeLimit,
		"mcp_config_path", cfg.MCP.ConfigPath,
		"mcp_server_count", len(cfg.MCP.Servers),
	)
	logger.DebugContext(ctx, "Starting agent app", "llm_provider", cfg.Provider)

	provider, err := newProvider(cfg, logger)
	if err != nil {
		return err
	}

	session := runtime.NewSession()
	memoryStore := newMemoryStore(cfg)
	memoryContext := memory.Context{}

	memoryContext, err = memory.NewLoader(memoryStore).Load(ctx)
	if err != nil {
		logger.DebugContext(ctx, "Memory context load failed", "error", err)
		return err
	}
	tools, mcpManager, closers, err := newToolSet(ctx, cfg, logger, memoryStore)
	if err != nil {
		logger.DebugContext(ctx, "Tool registry initialization failed", "error", err)
		return err
	}
	defer closeAll(ctx, logger, closers)
	logger.DebugContext(ctx, "Tool registry initialized", "tool_count", len(tools.Definitions()))
	promptBuilder := prompt.NewBuilder(newSessionPrompt(memoryContext, mcpManager.Catalog()))

	runtimeAgent := agent.NewReActAgent(
		provider,
		agent.WithPromptBuilder(promptBuilder),
		agent.WithTools(tools),
		agent.WithLogger(logger),
		agent.WithMaxRounds(cfg.Agent.MaxRounds),
		agent.WithToolTimeout(cfg.Agent.ToolTimeout),
	)

	runner := tui.NewRunner(session, runtimeAgent, stdin, stdout, logBuffer, initialPrompt)
	if err := runner.Run(ctx); err != nil {
		logger.DebugContext(ctx, "Agent TUI stopped with error", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent TUI stopped")
	return nil
}
