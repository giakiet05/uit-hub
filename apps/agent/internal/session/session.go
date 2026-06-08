// Package session owns state that lives for one chat session.
package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

// Snapshot is the persistable state of one session.
type Snapshot struct {
	ID               string
	StartedAt        time.Time
	Messages         []conversation.Message
	Usage            Usage
	RuntimeToolNames []string
}

// Store persists session snapshots. Storage backends such as SQLite,
// PostgreSQL, or MongoDB can implement this interface.
type Store interface {
	Save(ctx context.Context, snapshot Snapshot) error
}

// State groups one conversation with session-scoped prompt, tool, and usage
// state.
type State struct {
	mu              sync.RWMutex
	ID              string
	StartedAt       time.Time
	Conversation    conversation.Conversation
	PromptSnapshot  prompt.SystemPrompt
	Tools           *tool.ToolSet
	ResultBudgeter  *tool.ResultBudgeter
	Usage           Usage
	ConcurrentTools bool
}

// Usage records aggregate counters for one session.
type Usage struct {
	usage.TokenUsage
	Runs         int
	LLMCalls     int
	ToolCalls    int
	ToolFailures int
}

// AddRunStats accumulates one completed agent run into this usage value.
func (u *Usage) AddRunStats(stats agent.RunStats) {
	if u == nil {
		return
	}
	u.Runs++
	u.LLMCalls += stats.LLMCalls
	u.ToolCalls += stats.ToolCalls
	u.ToolFailures += stats.ToolFailures
	u.TokenUsage.Add(stats.TokenUsage)
}

// Option configures a session state at construction time.
type Option func(*State)

// WithPromptSnapshot sets the session-stable prompt snapshot.
func WithPromptSnapshot(promptSnapshot prompt.SystemPrompt) Option {
	return func(s *State) {
		s.PromptSnapshot = promptSnapshot
	}
}

// WithTools sets the session-visible tool collection.
func WithTools(tools *tool.ToolSet) Option {
	return func(s *State) {
		s.Tools = tools
	}
}

// WithResultBudgeter sets the session tool-result budgeting policy.
func WithResultBudgeter(budgeter *tool.ResultBudgeter) Option {
	return func(s *State) {
		s.ResultBudgeter = budgeter
	}
}

// WithConcurrentTools sets the concurrent tools flag.
func WithConcurrentTools(concurrent bool) Option {
	return func(s *State) {
		s.ConcurrentTools = concurrent
	}
}

// WithHistory initializes the session with loaded database state.
func WithHistory(id string, startedAt time.Time, messages []conversation.Message, usage Usage) Option {
	return func(s *State) {
		s.ID = id
		s.StartedAt = startedAt
		s.Usage = usage
		for _, msg := range messages {
			s.Conversation.Append(msg)
		}
	}
}

// NewState creates a session with a unique ID.
func NewState(opts ...Option) *State {
	b := make([]byte, 4)
	rand.Read(b)
	uniqueID := fmt.Sprintf("session-%s-%s", time.Now().Format("20060102150405"), hex.EncodeToString(b))

	state := &State{
		ID:             uniqueID,
		StartedAt:      time.Now().UTC(),
		Conversation:   conversation.NewConversation(),
		PromptSnapshot: prompt.SystemPrompt{},
	}
	for _, opt := range opts {
		opt(state)
	}
	return state
}

// ToolDefinitions returns all tools visible in this session.
func (s *State) ToolDefinitions() []tool.Definition {
	if s == nil || s.Tools == nil {
		return nil
	}
	return s.Tools.Definitions()
}

// ToolDefinition returns one visible tool definition by name.
func (s *State) ToolDefinition(name string) (tool.Definition, bool) {
	if s == nil || s.Tools == nil {
		return tool.Definition{}, false
	}
	return s.Tools.Definition(name)
}

// AddRunStats accumulates the metrics from a completed agent run.
func (s *State) AddRunStats(stats agent.RunStats) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Usage.AddRunStats(stats)
}

// GetUsage returns a safe copy of the current session usage metrics.
func (s *State) GetUsage() Usage {
	if s == nil {
		return Usage{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Usage
}

// Snapshot returns a copy of the session state suitable for persistence.
func (s *State) Snapshot() Snapshot {
	if s == nil {
		return Snapshot{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Snapshot{
		ID:               s.ID,
		StartedAt:        s.StartedAt,
		Messages:         s.Conversation.Messages(),
		Usage:            s.Usage,
		RuntimeToolNames: s.Tools.RuntimeToolNames(),
	}
}
