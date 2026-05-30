package agent

import "github.com/giakiet05/uit-hub/apps/agent/internal/conversation"

// AssistantToolCalls extracts tool calls from an assistant message.
func AssistantToolCalls(message conversation.Message) []conversation.ToolCall {
	assistant, ok := message.(conversation.AssistantMessage)
	if !ok {
		return nil
	}
	return assistant.ToolCalls
}

// ToolCallNames returns tool names in call order for compact logging.
func ToolCallNames(calls []conversation.ToolCall) []string {
	names := make([]string, 0, len(calls))
	for _, call := range calls {
		names = append(names, call.Name)
	}
	return names
}
