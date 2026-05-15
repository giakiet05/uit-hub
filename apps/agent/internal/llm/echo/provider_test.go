package echo

import (
	"context"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

func TestGenerateReturnsLastUserMessage(t *testing.T) {
	provider := NewProvider()

	response, err := provider.Generate(context.Background(), llm.GenerateRequest{
		Messages: []llm.Message{
			llm.NewTextMessage(llm.RoleUser, "first"),
			llm.NewTextMessage(llm.RoleAssistant, "ignored"),
			llm.NewTextMessage(llm.RoleUser, "second"),
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if got, want := response.Message.ContentText(), "received: second"; got != want {
		t.Fatalf("response text = %q, want %q", got, want)
	}
}
