package react

import (
	"context"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

// decisionSummary returns the assistant's visible tool-use summary for debug
// logging.
func decisionSummary(message conversation.Message) string {
	text := strings.TrimSpace(conversation.Text(message))
	if text == "" {
		return "empty"
	}
	return loop.PreviewText(text)
}

// trace writes a debug event with a stable event key.
func (a *ReActAgent) trace(ctx context.Context, event string, message string, attrs ...any) {
	attrs = append([]any{"event", event}, attrs...)
	a.logger.DebugContext(ctx, message, attrs...)
}

// logRunStats writes the summary counters for one completed or failed run.
func (a *ReActAgent) logRunStats(ctx context.Context, sessionID string, stats *agent.RunStats) {
	a.logger.InfoContext(
		ctx,
		"Agent run stats",
		"event", "react.run.stats",
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
