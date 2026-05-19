package prompt

import (
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func TestBuildAgentMessagesPrependsSystemMessage(t *testing.T) {
	builder := NewBuilder(SystemPrompt{})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	messages := builder.BuildAgentMessages(AgentTypeReAct, history)
	if got, want := len(messages), 2; got != want {
		t.Fatalf("messages len = %d, want %d", got, want)
	}
	if _, ok := messages[0].(conversation.SystemMessage); !ok {
		t.Fatalf("first message type = %T, want conversation.SystemMessage", messages[0])
	}
}

func TestBuildAgentMessagesIncludesReActPrompt(t *testing.T) {
	builder := NewBuilder(SystemPrompt{})

	messages := builder.BuildAgentMessages(AgentTypeReAct, nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "You run a tool-calling ReAct loop.") {
		t.Fatalf("system prompt missing ReAct instruction: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesSessionPromptAfterBoundary(t *testing.T) {
	builder := NewBuilder(SystemPrompt{
		SessionText: "You are a UIT Hub agent.",
	})

	messages := builder.BuildAgentMessages(AgentTypeReAct, nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, sessionContextHeader+"\n\nYou are a UIT Hub agent.") {
		t.Fatalf("system prompt missing session context boundary: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesUncachedPromptAfterBoundary(t *testing.T) {
	builder := NewBuilder(SystemPrompt{
		Uncached: []UncachedSystemPrompt{
			{
				Text:   "MCP tool definitions",
				Reason: "MCP tool definitions are user-specific.",
			},
		},
	})

	messages := builder.BuildAgentMessages(AgentTypeReAct, nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, uncachedContextHeader+"\n\nMCP tool definitions") {
		t.Fatalf("system prompt missing uncached context boundary: %q", systemText)
	}
	if strings.Contains(systemText, "user-specific") {
		t.Fatalf("uncached reason leaked into system prompt: %q", systemText)
	}
}

func TestBuildAgentMessagesDoesNotMutateHistory(t *testing.T) {
	builder := NewBuilder(SystemPrompt{})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	_ = builder.BuildAgentMessages(AgentTypeReAct, history)

	if got, want := len(history), 1; got != want {
		t.Fatalf("history len = %d, want %d", got, want)
	}
	if got, want := conversation.Text(history[0]), "hello"; got != want {
		t.Fatalf("history text = %q, want %q", got, want)
	}
}
