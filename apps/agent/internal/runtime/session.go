// Package runtime owns per-run runtime objects such as sessions.
package runtime

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

var sessionCounter atomic.Uint64

// Session groups one conversation history with a stable runtime identifier.
type Session struct {
	ID           string
	StartedAt    time.Time
	Conversation conversation.Conversation
}

// NewSession creates a session with a monotonic process-local ID.
func NewSession() *Session {
	id := sessionCounter.Add(1)
	return &Session{
		ID:           "session-" + strconv.FormatUint(id, 10),
		StartedAt:    time.Now().UTC(),
		Conversation: conversation.NewConversation(),
	}
}
