package agent

import (
	"context"
	"log/slog"
)

// Trace writes a debug event with a stable dot-separated event key.
func Trace(logger *slog.Logger, ctx context.Context, event string, message string, attrs ...any) {
	attrs = append([]any{"event", event}, attrs...)
	logger.DebugContext(ctx, message, attrs...)
}

// LogRunStats writes aggregate counters for one completed or failed agent run.
// agentType is the dot-separated prefix used in the event key, e.g. "react" or
// "plan_execute".
func LogRunStats(logger *slog.Logger, agentType string, ctx context.Context, sessionID string, stats *RunStats) {
	logger.InfoContext(
		ctx,
		"Agent run stats",
		"event", agentType+".run.stats",
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
