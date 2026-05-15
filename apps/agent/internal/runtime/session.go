package runtime

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

var sessionCounter atomic.Uint64

type Session struct {
	ID        string
	StartedAt time.Time
	messages  []llm.Message
}

func NewSession() *Session {
	id := sessionCounter.Add(1)
	return &Session{
		ID:        "session-" + strconv.FormatUint(id, 10),
		StartedAt: time.Now().UTC(),
	}
}

func (s *Session) Append(message llm.Message) {
	s.messages = append(s.messages, message)
}

func (s *Session) Messages() []llm.Message {
	messages := make([]llm.Message, len(s.messages))
	copy(messages, s.messages)
	return messages
}
