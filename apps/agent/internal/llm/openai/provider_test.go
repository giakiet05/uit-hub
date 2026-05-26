package openai

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func TestToResponsesInputMapsMessages(t *testing.T) {
	input := toResponsesInput([]conversation.Message{
		conversation.NewSystemMessage("system prompt"),
		conversation.NewUserMessage("hello"),
	})

	if len(input) != 2 {
		t.Fatalf("len(input) = %d, want 2", len(input))
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if got := string(body); !strings.Contains(got, `"role":"system"`) {
		t.Fatalf("marshaled input = %s, want system role", got)
	}
	if got := string(body); !strings.Contains(got, `"content":"hello"`) {
		t.Fatalf("marshaled input = %s, want user content", got)
	}
}

func TestToResponsesRole(t *testing.T) {
	tests := []struct {
		name string
		role conversation.Role
		want string
	}{
		{name: "system", role: conversation.RoleSystem, want: "system"},
		{name: "assistant", role: conversation.RoleAssistant, want: "assistant"},
		{name: "user", role: conversation.RoleUser, want: "user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(toResponsesRole(tt.role)); got != tt.want {
				t.Fatalf("toResponsesRole(%q) = %q, want %q", tt.role, got, tt.want)
			}
		})
	}
}
