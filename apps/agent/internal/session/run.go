package session

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
)

// Session owns one chat session and delegates user prompts to an agent service.
type Session struct {
	state *State
	agent agent.Agent
	bus   *bus.EventBus
}

// NewSession creates a session runtime around session state and a stateless
// agent service.
func NewSession(state *State, runtimeAgent agent.Agent, bus *bus.EventBus) *Session {
	return &Session{
		state: state,
		agent: runtimeAgent,
		bus:   bus,
	}
}

// SetAgent dynamically replaces the stateless agent service for this session.
// It is safe to call before the next Run.
func (s *Session) SetAgent(a agent.Agent) {
	if s != nil {
		s.agent = a
	}
}

// Bus returns the EventBus attached to this session.
func (s *Session) Bus() *bus.EventBus {
	if s == nil {
		return nil
	}
	return s.bus
}

// State returns the session-owned state snapshot for read-only callers.
func (s *Session) State() *State {
	if s == nil {
		return nil
	}
	return s.state
}

// Run builds one agent run input from session state and forwards agent events.
func (s *Session) Run(ctx context.Context, userPrompt string) {
	if s == nil || s.state == nil || s.agent == nil || s.bus == nil {
		return
	}

	input := agent.RunInput{
		SessionID:       s.state.ID,
		UserPrompt:      userPrompt,
		Conversation:    &s.state.Conversation,
		PromptSnapshot:  s.state.PromptSnapshot,
		Tools:           s.state.Tools,
		ResultBudgeter:  s.state.ResultBudgeter,
		ConcurrentTools: s.state.ConcurrentTools,
	}

	for event := range s.agent.Run(ctx, input) {
		s.bus.Publish(event)
	}
}
