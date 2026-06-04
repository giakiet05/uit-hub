// Package app wires configuration, providers, tools, sessions, and the TUI.
package app

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/planexecute"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/react"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
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
		"llm_provider", os.Getenv("DEFAULT_PROVIDER"),
		"agent_model", os.Getenv("DEFAULT_MODEL"),
		"agent_max_rounds", cfg.Agent.MaxRounds,
		"agent_tool_timeout", cfg.Agent.ToolTimeout.String(),
		"memory_path", cfg.Memory.Path,
		"tool_runtime_limit", cfg.Tool.RuntimeLimit,
		"compaction_tool_result_budget_enabled", cfg.Compaction.ToolResultBudgetEnabled,
		"compaction_tool_result_budget_max_chars", cfg.Compaction.ToolResultBudgetMaxChars,
		"compaction_snip_enabled", cfg.Compaction.SnipEnabled,
		"compaction_snip_trigger_ratio", cfg.Compaction.SnipTriggerRatio,
		"compaction_microcompact_enabled", cfg.Compaction.MicrocompactEnabled,
		"compaction_microcompact_trigger_ratio", cfg.Compaction.MicrocompactTriggerRatio,
		"compaction_microcompact_min_chars", cfg.Compaction.MicrocompactMinChars,
		"compaction_microcompact_keep_recent_tool_results", cfg.Compaction.MicrocompactKeepRecentToolResults,
		"compaction_context_collapse_enabled", cfg.Compaction.ContextCollapseEnabled,
		"compaction_context_collapse_trigger_ratio", cfg.Compaction.ContextCollapseTriggerRatio,
		"compaction_context_collapse_target_ratio", cfg.Compaction.ContextCollapseTargetRatio,
		"compaction_context_collapse_keep_recent_turns", cfg.Compaction.ContextCollapseKeepRecentTurns,
		"compaction_context_collapse_max_segments", cfg.Compaction.ContextCollapseMaxSegments,
		"compaction_context_collapse_max_concurrency", cfg.Compaction.ContextCollapseMaxConcurrency,
		"compaction_context_collapse_min_segment_tokens", cfg.Compaction.ContextCollapseMinSegmentTokens,
		"compaction_auto_compact_enabled", cfg.Compaction.AutoCompactEnabled,
		"mcp_config_path", cfg.MCP.ConfigPath,
		"mcp_server_count", len(cfg.MCP.Servers),
	)
	logger.DebugContext(ctx, "Starting agent app", "llm_provider", os.Getenv("DEFAULT_PROVIDER"))

	runtime, err := NewRuntime(ctx, cfg, logger, resume, resumeID)
	if err != nil {
		return err
	}
	defer runtime.Close(ctx)

	providerName := os.Getenv("DEFAULT_PROVIDER")
	if providerName == "" {
		providerName = "openai"
	}
	modelName := os.Getenv("DEFAULT_MODEL")
	if modelName == "" {
		modelName = "gpt-5.4-mini"
	}

	model, err := cfg.Registry.GetModel(providerName, modelName)
	if err != nil {
		logger.DebugContext(ctx, "Failed to get model from registry", "error", err)
		return err
	}

	var runtimeAgent agent.Agent
	switch cfg.Agent.Type {
	case agent.TypePlanAndExecute:
		runtimeAgent = planexecute.NewPlanAndExecuteAgent(
			model,
			planexecute.WithPlanMaxSteps(cfg.Agent.MaxRounds),
			planexecute.WithPlanToolTimeout(cfg.Agent.ToolTimeout),
		)
	default:
		runtimeAgent = react.NewReActAgent(
			model,
			react.WithMaxRounds(cfg.Agent.MaxRounds),
			react.WithToolTimeout(cfg.Agent.ToolTimeout),
			react.WithCompaction(conversation.CompactionOptions{
				DisableSnip:                       !cfg.Compaction.SnipEnabled,
				DisableMicrocompact:               !cfg.Compaction.MicrocompactEnabled,
				MicrocompactMinChars:              cfg.Compaction.MicrocompactMinChars,
				MicrocompactKeepRecentToolResults: cfg.Compaction.MicrocompactKeepRecentToolResults,
				MicrocompactTriggerRatio:          cfg.Compaction.MicrocompactTriggerRatio,
				SnipTriggerRatio:                  cfg.Compaction.SnipTriggerRatio,
			}),
			react.WithContextCollapse(react.ContextCollapseOptions{
				Enabled:          cfg.Compaction.ContextCollapseEnabled,
				TriggerRatio:     cfg.Compaction.ContextCollapseTriggerRatio,
				TargetRatio:      cfg.Compaction.ContextCollapseTargetRatio,
				KeepRecentTurns:  cfg.Compaction.ContextCollapseKeepRecentTurns,
				MaxSegments:      cfg.Compaction.ContextCollapseMaxSegments,
				MaxConcurrency:   cfg.Compaction.ContextCollapseMaxConcurrency,
				MinSegmentTokens: cfg.Compaction.ContextCollapseMinSegmentTokens,
			}),
		)
	}

	runtime.Session.SetAgent(runtimeAgent)

	runner := tui.NewRunner(runtime.Session, stdin, stdout, logBuffer, initialPrompt)
	if err := runner.Run(ctx); err != nil {
		logger.DebugContext(ctx, "Agent TUI stopped with error", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent TUI stopped")
	return nil
}
