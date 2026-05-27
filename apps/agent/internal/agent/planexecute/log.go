package planexecute

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
)

func (a *PlanAndExecuteAgent) trace(ctx context.Context, event string, message string, attrs ...any) {
	attrs = append([]any{"event", event}, attrs...)
	a.logger.DebugContext(ctx, message, attrs...)
}

func (a *PlanAndExecuteAgent) logRunStats(ctx context.Context, sessionID string, stats *agent.RunStats) {
	a.logger.InfoContext(
		ctx,
		"Agent run stats",
		"event", "plan_execute.run.stats",
		"session_id", sessionID,
		"duration", stats.Duration().String(),
		"rounds", stats.Rounds,
		"llm_calls", stats.LLMCalls,
		"tool_calls", stats.ToolCalls,
		"tool_failures", stats.ToolFailures,
		"input_tokens", stats.InputTokens,
		"output_tokens", stats.OutputTokens,
	)
}
