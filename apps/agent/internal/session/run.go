package session

import (
	"context"
	"log/slog"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
)

// Session owns one chat session and delegates user prompts to an agent service.
type Session struct {
	state           *State
	agent           agent.Agent
	bus             *bus.EventBus
	store           Store
	logger          *slog.Logger
	shutdownContext context.Context
}

// RuntimeOption configures session runtime dependencies.
type RuntimeOption func(*Session)

// WithStore sets the persistent session store.
func WithStore(store Store) RuntimeOption {
	return func(s *Session) {
		s.store = store
	}
}

// WithLogger sets the session runtime logger.
func WithLogger(logger *slog.Logger) RuntimeOption {
	return func(s *Session) {
		s.logger = logger
	}
}

// WithShutdownContext sets the app-level context used for cleanup and
// force-finish side-effecting work.
func WithShutdownContext(ctx context.Context) RuntimeOption {
	return func(s *Session) {
		s.shutdownContext = ctx
	}
}

// NewSession creates a session runtime around session state and a stateless
// agent service.
func NewSession(state *State, runtimeAgent agent.Agent, eventBus *bus.EventBus, opts ...RuntimeOption) *Session {
	session := &Session{
		state: state,
		agent: runtimeAgent,
		bus:   eventBus,
	}
	for _, opt := range opts {
		opt(session)
	}
	if session.logger == nil {
		session.logger = slog.Default()
	}
	return session
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
		ShutdownContext: s.shutdownContext,
	}
	if input.ShutdownContext == nil {
		input.ShutdownContext = ctx
	}

	for event := range s.agent.Run(ctx, input) {
		s.handleRunEvent(ctx, event)
		s.bus.Publish(event)
	}
}

func (s *Session) handleRunEvent(ctx context.Context, event agent.Event) {
	switch typed := event.(type) {
	case agent.RunCompletedEvent:
		s.finishRun(ctx, typed.Stats)
	case agent.RunFailedEvent:
		s.finishRun(ctx, typed.Stats)
	}
}

func (s *Session) finishRun(ctx context.Context, stats agent.RunStats) {
	s.state.AddRunStats(stats)
	if s.store == nil {
		return
	}

	saveCtx := ctx
	if ctx.Err() != nil {
		var cancel context.CancelFunc
		saveCtx, cancel = context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
	}

	snapshot := s.state.Snapshot()
	if err := s.store.Save(saveCtx, snapshot); err != nil {
		s.logger.ErrorContext(context.WithoutCancel(ctx), "Failed to save session", "session_id", snapshot.ID, "error", err)
		return
	}
	s.logger.DebugContext(context.WithoutCancel(ctx), "Session saved", "session_id", snapshot.ID, "message_count", len(snapshot.Messages))
}
