package loop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// Emit sends one agent event unless the context is canceled.
func Emit(ctx context.Context, events chan<- agent.Event, event agent.Event) bool {
	select {
	case <-ctx.Done():
		return false
	case events <- event:
		return true
	}
}

// AssistantToolCalls extracts tool calls from an assistant message.
func AssistantToolCalls(message conversation.Message) []conversation.ToolCall {
	assistant, ok := message.(conversation.AssistantMessage)
	if !ok {
		return nil
	}
	return assistant.ToolCalls
}

// NormalizeToolResult fills missing result metadata from the original tool call.
func NormalizeToolResult(call conversation.ToolCall, result tool.Result) tool.Result {
	if result.CallID == "" {
		result.CallID = call.ID
	}
	if result.Name == "" {
		result.Name = call.Name
	}
	return result
}

// ToolErrorObservation converts internal tool failures into sanitized
// observations that can be safely shown back to the model.
func ToolErrorObservation(err error) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "tool error: execution timed out"
	case errors.Is(err, tool.ErrNotFound):
		return "tool error: requested tool is not available"
	default:
		return "tool error: execution failed"
	}
}

// ToolCallNames returns tool names in call order for compact logging.
func ToolCallNames(calls []conversation.ToolCall) []string {
	names := make([]string, 0, len(calls))
	for _, call := range calls {
		names = append(names, call.Name)
	}
	return names
}

// PreviewValue renders a structured value as a single-line log value.
func PreviewValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return PreviewText(fmt.Sprintf("%v", value))
	}
	return PreviewText(string(data))
}

// PreviewText returns a single-line value suitable for structured logs.
func PreviewText(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}
