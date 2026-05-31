package storage_test

import (
	"os"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

func TestUpsertSession(t *testing.T) {
	testDB := "test_agent.db"
	defer os.Remove(testDB)

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

	err = db.UpsertSession("session-1", startedAt, messages, usage)
	if err != nil {
		t.Fatalf("Failed to upsert session: %v", err)
	}
}
