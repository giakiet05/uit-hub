package prompt_test

import (
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
)

func TestBuildAgentMessagesPrependsSystemMessage(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SessionPrompt{
		StaticParts: []prompt.StaticPart{"system"},
	})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	messages := builder.BuildAgentMessages(nil, history)
	if got, want := len(messages), 2; got != want {
		t.Fatalf("messages len = %d, want %d", got, want)
	}
	if _, ok := messages[0].(conversation.SystemMessage); !ok {
		t.Fatalf("first message type = %T, want conversation.SystemMessage", messages[0])
	}
}

func TestBuildAgentMessagesIncludesReActPrompt(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SessionPrompt{
		StaticParts: []prompt.StaticPart{prompt.ReActStaticPrompt()},
	})

	messages := builder.BuildAgentMessages(nil, nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "You run a tool-calling ReAct loop.") {
		t.Fatalf("system prompt missing ReAct instruction: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesSessionPromptAfterBoundary(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SessionPrompt{
		DynamicParts: []prompt.DynamicPart{"You are a UIT Hub agent."},
	})

	messages := builder.BuildAgentMessages(nil, nil)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "=== DYNAMIC CONTEXT ===\n\nYou are a UIT Hub agent.") {
		t.Fatalf("system prompt missing session context boundary: %q", systemText)
	}
}

func TestBuildAgentMessagesIncludesUncachedPromptAfterBoundary(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SessionPrompt{})

	messages := builder.BuildAgentMessages(
		[]prompt.UncachedPart{"MCP tool definitions"},
		nil,
	)
	systemText := conversation.Text(messages[0])
	if !strings.Contains(systemText, "=== UNCACHED CONTEXT ===\n\nMCP tool definitions") {
		t.Fatalf("system prompt missing uncached context boundary: %q", systemText)
	}
}

func TestBuildAgentMessagesDoesNotMutateHistory(t *testing.T) {
	builder := prompt.NewBuilder(prompt.SessionPrompt{})
	history := []conversation.Message{
		conversation.NewUserMessage("hello"),
	}

	_ = builder.BuildAgentMessages(nil, history)

	if got, want := len(history), 1; got != want {
		t.Fatalf("history len = %d, want %d", got, want)
	}
	if got, want := conversation.Text(history[0]), "hello"; got != want {
		t.Fatalf("history text = %q, want %q", got, want)
	}
}

func TestReActStaticPrompt(t *testing.T) {
	staticPrompt := prompt.ReActStaticPrompt()

	if !strings.Contains(string(staticPrompt), "You run a tool-calling ReAct loop.") {
		t.Fatalf("ReActStaticPrompt() missing ReAct instruction: %q", staticPrompt)
	}
}

func TestMemoryDynamicPromptIncludesIndexAndPolicy(t *testing.T) {
	dynamicPrompt := prompt.MemoryDynamicPrompt(memory.Context{
		IndexText: "- [Response Style](feedback_response_style) [feedback] -- User prefers concise answers.",
	})

	if !strings.Contains(string(dynamicPrompt), "## Long-Term Memory") {
		t.Fatalf("memory prompt missing header: %q", dynamicPrompt)
	}
	if !strings.Contains(string(dynamicPrompt), "feedback_response_style") {
		t.Fatalf("memory prompt missing index: %q", dynamicPrompt)
	}
	if !strings.Contains(string(dynamicPrompt), "Do not save transient task data, student records") {
		t.Fatalf("memory prompt missing save policy: %q", dynamicPrompt)
	}
}

func TestMemoryDynamicPromptSkipsEmptyContext(t *testing.T) {
	dynamicPrompt := prompt.MemoryDynamicPrompt(memory.Context{})

	if dynamicPrompt != "" {
		t.Fatalf("MemoryDynamicPrompt() = %q, want empty", dynamicPrompt)
	}
}

func TestMCPToolCatalogDynamicPrompt(t *testing.T) {
	dynamicPrompt := prompt.MCPToolCatalogDynamicPrompt([]mcpadapter.ToolMetadata{
		{
			Name:              "uit__get_student",
			ServerName:        "uit",
			ServerDescription: "UIT academic server.",
			Description:       "Read student profile.",
		},
	})

	text := string(dynamicPrompt)
	for _, want := range []string{
		"## MCP Tool Catalog",
		"### uit",
		"UIT academic server.",
		"Use load_mcp_tool",
		"uit__get_student: Read student profile.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("catalog prompt missing %q in %q", want, text)
		}
	}
}

func TestMCPToolCatalogDynamicPromptSkipsEmptyCatalog(t *testing.T) {
	dynamicPrompt := prompt.MCPToolCatalogDynamicPrompt(nil)

	if dynamicPrompt != "" {
		t.Fatalf("MCPToolCatalogDynamicPrompt() = %q, want empty", dynamicPrompt)
	}
}
