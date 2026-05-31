package react

import (
	"context"
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
	results, ok := agent.RunToolBatch(ctx, "react", events, input, round, calls, a.toolTimeout, stats)
	if !ok {
		return false
	}

	for _, result := range results {
		input.Conversation.Append(conversation.NewToolResultMessage(result.CallID, result.Content))
	}

	return true
}
