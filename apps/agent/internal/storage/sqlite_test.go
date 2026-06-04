package storage_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

func TestUpsertSession(t *testing.T) {
	testDB := filepath.Join(t.TempDir(), "test_agent.db")

	db, err := storage.InitDB(testDB)
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	startedAt := time.Now()
	messages := []conversation.Message{
		conversation.NewUserMessage("Hello"),
	}
	usage := session.Usage{}
	usage.Runs = 1

	err = db.UpsertSession(context.Background(), "session-1", startedAt, messages, usage, nil)
	if err != nil {
		t.Fatalf("Failed to upsert session: %v", err)
	}
}

func TestSessionRecordDecodeRestoresConversationAndUsage(t *testing.T) {
	testDB := filepath.Join(t.TempDir(), "test_agent.db")

	db, err := storage.InitDB(testDB)
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	startedAt := time.Now()
	messages := []conversation.Message{
		conversation.NewSystemMessage("system prompt should not be stored"),
		conversation.NewUserMessage("find student"),
		conversation.NewAssistantMessage("calling tool", []conversation.ToolCall{
			{
				ID:   "call-1",
				Name: "mock_uit__get_student_profile",
				Arguments: map[string]any{
					"student_id": "22520001",
				},
			},
		}),
		conversation.NewToolResultMessage("call-1", `{"student_id":"22520001","gpa":3.42}`),
		conversation.NewAssistantMessage("student GPA is 3.42", nil),
	}
	usage := session.Usage{}
	usage.Runs = 2
	usage.LLMCalls = 3
	usage.InputTokens = 100
	usage.OutputTokens = 20

	runtimeToolNames := []string{"mock_uit__get_student_profile"}
	if err := db.UpsertSession(context.Background(), "session-1", startedAt, messages, usage, runtimeToolNames); err != nil {
		t.Fatalf("Failed to upsert session: %v", err)
	}

	record, err := db.GetSession("session-1")
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	decodedMessages, decodedUsage, decodedRuntimeToolNames, err := record.Decode()
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if got, want := len(decodedMessages), 4; got != want {
		t.Fatalf("len(decodedMessages) = %d, want %d", got, want)
	}
	if _, ok := decodedMessages[0].(conversation.UserMessage); !ok {
		t.Fatalf("decodedMessages[0] = %T, want UserMessage", decodedMessages[0])
	}
	assistantMessage, ok := decodedMessages[1].(conversation.AssistantMessage)
	if !ok {
		t.Fatalf("decodedMessages[1] = %T, want AssistantMessage", decodedMessages[1])
	}
	if got, want := len(assistantMessage.ToolCalls), 1; got != want {
		t.Fatalf("len(ToolCalls) = %d, want %d", got, want)
	}
	if _, ok := decodedMessages[2].(conversation.ToolResultMessage); !ok {
		t.Fatalf("decodedMessages[2] = %T, want ToolResultMessage", decodedMessages[2])
	}
	if decodedUsage.Runs != usage.Runs || decodedUsage.LLMCalls != usage.LLMCalls {
		t.Fatalf("decodedUsage = %+v, want %+v", decodedUsage, usage)
	}
	if decodedUsage.InputTokens != usage.InputTokens || decodedUsage.OutputTokens != usage.OutputTokens {
		t.Fatalf("decoded token usage = %+v, want %+v", decodedUsage.TokenUsage, usage.TokenUsage)
	}
	if got, want := len(decodedRuntimeToolNames), len(runtimeToolNames); got != want {
		t.Fatalf("len(decodedRuntimeToolNames) = %d, want %d", got, want)
	}
	if decodedRuntimeToolNames[0] != runtimeToolNames[0] {
		t.Fatalf("decodedRuntimeToolNames[0] = %q, want %q", decodedRuntimeToolNames[0], runtimeToolNames[0])
	}
}
