package session

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
)

// Session owns one chat session and delegates user prompts to an agent service.
type Session struct {
	state *State
	agent agent.Agent
}

// NewSession creates a session runtime around session state and a stateless
// agent service.
func NewSession(state *State, runtimeAgent agent.Agent) *Session {
	return &Session{
		state: state,
		agent: runtimeAgent,
	}
}

// State returns the session-owned state snapshot for read-only callers.
func (s *Session) State() *State {
	if s == nil {
		return nil
	}
	return s.state
}

// Run builds one agent run input from session state and forwards agent events.
func (s *Session) Run(ctx context.Context, userPrompt string) <-chan agent.Event {
	events := make(chan agent.Event)
	go func() {
		defer close(events)
		if s == nil || s.state == nil || s.agent == nil {
			return
		}

		input := agent.RunInput{
			SessionID:       s.state.ID,
			UserPrompt:      userPrompt,
			Conversation:    &s.state.Conversation,
			PromptSnapshot:  s.state.PromptSnapshot,
			Tools:           s.state.Tools,
			ConcurrentTools: s.state.ConcurrentTools,
		}

		for event := range s.agent.Run(ctx, input) {
			s.recordUsage(event)
			select {
			case events <- event:
			case <-ctx.Done():
				return
			}
		}
	}()
	return events
}

func (s *Session) recordUsage(event agent.Event) {
	switch typed := event.(type) {
	case agent.RunCompletedEvent:
		s.addRunStats(typed.Stats)
	case agent.RunFailedEvent:
		s.addRunStats(typed.Stats)
	}
}

func (s *Session) addRunStats(stats agent.RunStats) {
	if s == nil || s.state == nil {
		return
	}
	s.state.Usage.Runs++
	s.state.Usage.LLMCalls += stats.LLMCalls
	s.state.Usage.ToolCalls += stats.ToolCalls
	s.state.Usage.ToolFailures += stats.ToolFailures
	s.state.Usage.TokenUsage.Add(stats.TokenUsage)
}
