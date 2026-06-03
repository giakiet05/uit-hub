package conversation

import (
	"strings"
	"testing"
)

func TestMicrocompactSkipsToolResultBelowMinChars(t *testing.T) {
	messages := []Message{
		NewUserMessage("question"),
		NewAssistantMessage("", []ToolCall{{ID: "call-1", Name: "rag_search"}}),
		NewToolResultMessage("call-1", strings.Repeat("a", 200)),
	}

	compacted := Microcompact(messages, 1, 1000, 0)

	if got := Text(compacted[2]); got != strings.Repeat("a", 200) {
		t.Fatalf("tool result = %q, want original content", got)
	}
}

func TestMicrocompactHidesToolResultAboveMinChars(t *testing.T) {
	messages := []Message{
		NewUserMessage("question"),
		NewAssistantMessage("", []ToolCall{{ID: "call-1", Name: "rag_search"}}),
		NewToolResultMessage("call-1", strings.Repeat("a", 200)),
	}

	compacted := Microcompact(messages, 1, 100, 0)

	got := Text(compacted[2])
	if !strings.Contains(got, "Tool result hidden") {
		t.Fatalf("tool result = %q, want compact notice", got)
	}
	if !strings.Contains(got, "call-1") {
		t.Fatalf("tool result = %q, want call id in compact notice", got)
	}
}

func TestMicrocompactKeepsRecentToolResults(t *testing.T) {
	messages := []Message{
		NewUserMessage("question"),
		NewAssistantMessage("", []ToolCall{{ID: "call-old", Name: "rag_search"}}),
		NewToolResultMessage("call-old", strings.Repeat("old", 100)),
		NewAssistantMessage("", []ToolCall{{ID: "call-new", Name: "rag_search"}}),
		NewToolResultMessage("call-new", strings.Repeat("new", 100)),
	}

	compacted := Microcompact(messages, 1, 100, 1)

	if got := Text(compacted[2]); !strings.Contains(got, "Tool result hidden") {
		t.Fatalf("old tool result = %q, want compact notice", got)
	}
	if got := Text(compacted[4]); got != strings.Repeat("new", 100) {
		t.Fatalf("recent tool result = %q, want original content", got)
	}
}

func TestSnipCompactRemovesWholeOldUserTurn(t *testing.T) {
	messages := []Message{
		NewUserMessage(strings.Repeat("old user ", 500)),
		NewAssistantMessage("old assistant will call tool", []ToolCall{{ID: "call-old", Name: "search"}}),
		NewToolResultMessage("call-old", strings.Repeat("old tool result ", 500)),
		NewAssistantMessage("old final answer", nil),
		NewUserMessage("latest user"),
		NewAssistantMessage("latest assistant will call tool", []ToolCall{{ID: "call-new", Name: "lookup"}}),
		NewToolResultMessage("call-new", "latest tool result"),
		NewAssistantMessage("latest final answer", nil),
	}

	compacted := SnipCompact(messages, 10)
	transcript := transcriptText(compacted)

	if strings.Contains(transcript, "call-old") || strings.Contains(transcript, "old tool result") {
		t.Fatalf("old tool trajectory was not fully removed:\n%s", transcript)
	}
	if !strings.Contains(transcript, "latest user") {
		t.Fatalf("latest user turn was removed:\n%s", transcript)
	}
	if !strings.Contains(transcript, "latest tool result") {
		t.Fatalf("latest tool result was removed:\n%s", transcript)
	}
}

func TestSnipCompactDoesNotRemoveOnlyUserTurn(t *testing.T) {
	messages := []Message{
		NewUserMessage(strings.Repeat("only user ", 500)),
		NewAssistantMessage("only assistant", nil),
	}

	compacted := SnipCompact(messages, 1)

	if len(compacted) != len(messages) {
		t.Fatalf("len(compacted) = %d, want %d", len(compacted), len(messages))
	}
}

func transcriptText(messages []Message) string {
	var builder strings.Builder
	for _, message := range messages {
		builder.WriteString(Text(message))
		if assistant, ok := message.(AssistantMessage); ok {
			for _, call := range assistant.ToolCalls {
				builder.WriteString(" ")
				builder.WriteString(call.ID)
				builder.WriteString(" ")
				builder.WriteString(call.Name)
			}
		}
		if toolResult, ok := message.(ToolResultMessage); ok {
			builder.WriteString(" ")
			builder.WriteString(toolResult.ToolCallID)
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
