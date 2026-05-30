package tui

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

func TestModelRendersConversationAndLogs(t *testing.T) {
	logs := NewLogBuffer(10)
	_, err := logs.Write([]byte("first log line\nsecond log line\n"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	m := newModel(context.Background(), newTestSession(), logs, "")
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
		"last i/o 0/0 | session i/o 0/0",
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

func TestModelRendersAgentActivityEvents(t *testing.T) {
	m := newModel(context.Background(), newTestSession(), nil, "")
	m.width = 80
	m.height = 24
	m.resizeViewports()

	m = m.handleAgentEvent(agent.ModelCallStartedEvent{SessionID: "session-1", Round: 1})
	m = m.handleAgentEvent(agent.ModelTextDeltaEvent{
		SessionID: "session-1",
		Round:     1,
		Delta:     "Đang tính",
	})
	m = m.handleAgentEvent(agent.ModelCallCompletedEvent{
		SessionID:     "session-1",
		Round:         1,
		ToolCallNames: []string{"calculator"},
		MessageText:   "Mình sẽ dùng calculator để tính.",
	})
	m = m.handleAgentEvent(agent.ToolCallStartedEvent{
		SessionID: "session-1",
		Round:     1,
		Call: conversation.ToolCall{
			ID:   "call-1",
			Name: "calculator",
			Arguments: map[string]any{
				"operation": "add",
				"a":         1,
				"b":         2,
			},
		},
	})
	m = m.handleAgentEvent(agent.RunCompletedEvent{
		SessionID: "session-1",
		Reason:    agent.TerminalCompleted,
		Stats: agent.RunStats{
			TokenUsage: usage.TokenUsage{
				InputTokens:  12,
				OutputTokens: 7,
			},
		},
	})

	view := m.View()
	expectedParts := []string{
		"- thinking",
		"Đang tính",
		"model response: Mình sẽ dùng",
		"calculator để tính.",
		"- tool calculator",
		`{"a":1,"b":2,"operation":"add"}`,
		"last i/o 12/7 | session i/o 12/7",
	}
	for _, part := range expectedParts {
		if !strings.Contains(view, part) {
			t.Fatalf("view missing %q in:\n%s", part, view)
		}
	}
}

type fakeAgent struct{}

func (fakeAgent) Run(ctx context.Context, input agent.RunInput) <-chan agent.Event {
	events := make(chan agent.Event, 2)
	go func() {
		defer close(events)
		events <- agent.FinalAnswerEvent{
			SessionID: input.SessionID,
			Round:     1,
			Message:   conversation.NewAssistantMessage("ok", nil),
		}
		events <- agent.RunCompletedEvent{
			SessionID: input.SessionID,
			Reason:    agent.TerminalCompleted,
			Stats: agent.RunStats{
				TokenUsage: usage.TokenUsage{
					InputTokens:  3,
					OutputTokens: 2,
				},
			},
		}
	}()
	return events
}

func newTestSession() *session.Session {
	return session.NewSession(session.NewState(), fakeAgent{})
}
