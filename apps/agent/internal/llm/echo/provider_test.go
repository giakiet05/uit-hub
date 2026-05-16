package echo

import (
	"context"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

func TestGenerateReturnsLastUserMessage(t *testing.T) {
	provider := NewProvider()

	response, err := provider.Generate(context.Background(), llm.GenerateRequest{
		Messages: []conversation.Message{
			conversation.NewUserMessage("first"),
			conversation.NewAssistantMessage("ignored", nil),
			conversation.NewUserMessage("second"),
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if got, want := conversation.Text(response.Message), "received: second"; got != want {
		t.Fatalf("response text = %q, want %q", got, want)
	}
}
