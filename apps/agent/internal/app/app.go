// Package app wires configuration, providers, tools, sessions, and the TUI.
package app

import (
	"context"
	"io"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tui"
)

// Run starts the agent application in TUI mode.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	var resume bool
	var resumeID string
	if len(args) > 0 && args[0] == "resume" {
		resume = true
		args = args[1:]
		if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
			resumeID = args[0]
			args = args[1:]
		} else {
			// They want to resume but didn't provide an ID, open picker.
			db, err := storage.InitDB("agent.db")
			if err == nil {
				sessions, err := db.ListSessions(10) // show up to 10 recent
				if err == nil && len(sessions) > 0 {
					resumeID = selectSessionTUI(sessions)
					if resumeID == "" {
						// User canceled picker
						return nil
					}
				}
			}
		}
	}

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

	runtime, err := NewRuntime(ctx, cfg, logger, resume, resumeID)
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
