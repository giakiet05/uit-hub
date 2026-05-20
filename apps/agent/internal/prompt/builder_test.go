package prompt_test

import (
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
)

func TestBuildAgentMessagesPrependsSystemMessage(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SystemPrompt{})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	messages := builder.BuildAgentMessages(agent.TypeReAct.String(), history)
	if got, want := len(messages), 2; got != want {
		t.Fatalf("messages len = %d, want %d", got, want)
	}
	if _, ok := messages[0].(conversation.SystemMessage); !ok {
		t.Fatalf("first message type = %T, want conversation.SystemMessage", messages[0])
	}
}

func TestBuildAgentMessagesIncludesReActPrompt(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SystemPrompt{})

	messages := builder.BuildAgentMessages(agent.TypeReAct.String(), nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "You run a tool-calling ReAct loop.") {
		t.Fatalf("system prompt missing ReAct instruction: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesSessionPromptAfterBoundary(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SystemPrompt{
		SessionText: "You are a UIT Hub agent.",
	})

	messages := builder.BuildAgentMessages(agent.TypeReAct.String(), nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "=== SESSION CONTEXT ===\n\nYou are a UIT Hub agent.") {
		t.Fatalf("system prompt missing session context boundary: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesUncachedPromptAfterBoundary(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SystemPrompt{
		Uncached: []prompt.UncachedSystemPrompt{
			{
				Text:   "MCP tool definitions",
				Reason: "MCP tool definitions are user-specific.",
			},
		},
	})

	messages := builder.BuildAgentMessages(agent.TypeReAct.String(), nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "=== UNCACHED CONTEXT ===\n\nMCP tool definitions") {
		t.Fatalf("system prompt missing uncached context boundary: %q", systemText)
	}
	if strings.Contains(systemText, "user-specific") {
		t.Fatalf("uncached reason leaked into system prompt: %q", systemText)
	}
}

func TestBuildAgentMessagesDoesNotMutateHistory(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SystemPrompt{})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	_ = builder.BuildAgentMessages(agent.TypeReAct.String(), history)

	if got, want := len(history), 1; got != want {
		t.Fatalf("history len = %d, want %d", got, want)
	}
	if got, want := conversation.Text(history[0]), "hello"; got != want {
		t.Fatalf("history text = %q, want %q", got, want)
	}
}

func TestDefaultSystemPromptUsesSessionText(t *testing.T) {
	systemPrompt := prompt.DefaultSystemPrompt()

	if systemPrompt.StableText != "" {
		t.Fatalf("StableText = %q, want empty", systemPrompt.StableText)
	}
	if systemPrompt.SessionText == "" {
		t.Fatal("SessionText is empty")
	}
	if len(systemPrompt.Uncached) != 0 {
		t.Fatalf("Uncached len = %d, want 0", len(systemPrompt.Uncached))
	}
}
