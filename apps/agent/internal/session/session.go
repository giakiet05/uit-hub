// Package session owns state that lives for one chat session.
package session

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

var sessionCounter atomic.Uint64

// State groups one conversation with session-scoped prompt, tool, and usage
// state.
type State struct {
	ID             string
	StartedAt      time.Time
	Conversation   conversation.Conversation
	PromptSnapshot prompt.SessionPrompt
	Tools          *tool.ToolSet
	Usage          Usage
}

// Usage records aggregate counters for one session.
type Usage struct {
	usage.TokenUsage
	Runs         int
	LLMCalls     int
	ToolCalls    int
	ToolFailures int
}

// UsageDelta records one run's contribution to session usage.
type UsageDelta struct {
	usage.TokenUsage
	LLMCalls     int
	ToolCalls    int
	ToolFailures int
}

// Option configures a session state at construction time.
type Option func(*State)

// WithPromptSnapshot sets the session-stable prompt snapshot.
func WithPromptSnapshot(promptSnapshot prompt.SessionPrompt) Option {
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

// NewState creates a session with a monotonic process-local ID.
func NewState(opts ...Option) *State {
	id := sessionCounter.Add(1)
	state := &State{
		ID:             "session-" + strconv.FormatUint(id, 10),
		StartedAt:      time.Now().UTC(),
		Conversation:   conversation.NewConversation(),
		PromptSnapshot: prompt.SessionPrompt{},
	}
	for _, opt := range opts {
		opt(state)
	}
	return state
}

// BuildAgentMessages prepends the session prompt snapshot to the conversation
// history for one LLM call.
func (s *State) BuildAgentMessages(uncached []prompt.UncachedPart) []conversation.Message {
	if s == nil {
		return nil
	}
	return prompt.NewBuilder(s.PromptSnapshot).BuildAgentMessages(uncached, s.Conversation.Messages())
}

// BuildMessages prepends the session prompt snapshot to custom messages for
// internal agent calls such as planning and finalization.
func (s *State) BuildMessages(uncached []prompt.UncachedPart, messages []conversation.Message) []conversation.Message {
	if s == nil {
		return nil
	}
	return prompt.NewBuilder(s.PromptSnapshot).BuildAgentMessages(uncached, messages)
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

// ExecuteTool runs a visible tool in this session.
func (s *State) ExecuteTool(ctx context.Context, call tool.Call) (tool.Result, error) {
	if s == nil || s.Tools == nil {
		return tool.Result{}, tool.ErrNotFound
	}
	return s.Tools.Execute(ctx, call)
}

// MarkToolUsed records usage of a session runtime tool for LRU eviction.
func (s *State) MarkToolUsed(name string) {
	if s == nil || s.Tools == nil {
		return
	}
	s.Tools.MarkUsed(name)
}

// AddUsage merges one run's usage into the session aggregate.
func (s *State) AddUsage(delta UsageDelta) {
	if s == nil {
		return
	}
	s.Usage.Runs++
	s.Usage.LLMCalls += delta.LLMCalls
	s.Usage.ToolCalls += delta.ToolCalls
	s.Usage.ToolFailures += delta.ToolFailures
	s.Usage.TokenUsage.Add(delta.TokenUsage)
}
