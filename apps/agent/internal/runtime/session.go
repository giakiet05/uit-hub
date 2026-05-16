package runtime

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

var sessionCounter atomic.Uint64

type Session struct {
	ID           string
	StartedAt    time.Time
	Conversation conversation.Conversation
}

func NewSession() *Session {
	id := sessionCounter.Add(1)
	return &Session{
		ID:           "session-" + strconv.FormatUint(id, 10),
		StartedAt:    time.Now().UTC(),
		Conversation: conversation.NewConversation(),
	}
}
