package tui

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

func TestModelRendersConversationAndLogs(t *testing.T) {
	logs := NewLogBuffer(10)
	_, err := logs.Write([]byte("first log line\nsecond log line\n"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	m := newModel(context.Background(), runtime.NewSession(), fakeAgent{}, logs, "")
	m.width = 80
	m.height = 24
	m.resizeViewports()
	m.conversation = append(
		m.conversation,
		conversationItem{role: conversationRoleUser, text: "hello"},
		conversationItem{role: conversationRoleAssistant, text: "xin chao"},
	)
	m.syncConversation(true)
	m.logLines = logs.Lines()
	m.syncLogs(true)

	view := m.View()
	expectedParts := []string{
		"Conversation",
		"> hello",
		"xin chao",
		"Logs",
		"first log line",
	}
	for _, part := range expectedParts {
		if !strings.Contains(view, part) {
			t.Fatalf("view missing %q in:\n%s", part, view)
		}
	}
	if strings.Contains(view, "user:") || strings.Contains(view, "assistant:") {
		t.Fatalf("view should not render role prefixes:\n%s", view)
	}
}

type fakeAgent struct{}

func (fakeAgent) Run(ctx context.Context, session *runtime.Session, input string) (conversation.Message, error) {
	return conversation.NewAssistantMessage("ok", nil), nil
}
