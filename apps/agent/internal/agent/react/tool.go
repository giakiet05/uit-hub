package react

import (
	"context"
	"fmt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

// executeToolCalls runs every requested tool call and appends each observation
// back into the session conversation.
func (a *ReActAgent) executeToolCalls(
	ctx context.Context,
	events chan<- agent.Event,
	input agent.RunInput,
	round int,
	calls []conversation.ToolCall,
	stats *agent.RunStats,
) bool {
	results, ok := agent.RunToolBatch(ctx, a.logger, "react", events, input, round, calls, a.toolTimeout, stats)
	if !ok {
		return false
	}

	for _, result := range results {
		input.Conversation.Append(conversation.NewToolResultMessage(result.CallID, result.Content))
		agent.Trace(
			a.logger,
			ctx,
			"react.tool.observation_appended",
			fmt.Sprintf("Round %d: Observe %s", round, result.Name),
			"session_id", input.SessionID,
			"round", round,
			"tool_call_id", result.CallID,
			"tool_name", result.Name,
			"result_chars", len(result.Content),
			"result_preview", agent.PreviewText(result.Content),
			// toolDuration is no longer available here per result, but that's okay for trace.
		)
	}

	return true
}
