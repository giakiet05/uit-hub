package session

import (
	"context"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
)

func TestSessionSavesAfterCanceledRunContext(t *testing.T) {
	runCtx, cancel := context.WithCancel(context.Background())
	cancel()

	store := &captureStore{}
	session := NewSession(
		NewState(),
		completedAgent{},
		bus.New(8, nil),
		WithStore(store),
	)

	session.Run(runCtx, "hello")

	if !store.saved {
		t.Fatal("store was not called")
	}
	if store.saveContextErr != nil {
		t.Fatalf("save context error = %v, want nil", store.saveContextErr)
	}
}

type captureStore struct {
	saved          bool
	saveContextErr error
}

func (s *captureStore) Save(ctx context.Context, snapshot Snapshot) error {
	s.saved = true
	s.saveContextErr = ctx.Err()
	return nil
}

type completedAgent struct{}

func (completedAgent) Run(ctx context.Context, input agent.RunInput) <-chan agent.Event {
	events := make(chan agent.Event, 1)
	go func() {
		defer close(events)
		events <- agent.RunCompletedEvent{
			SessionID: input.SessionID,
			Reason:    agent.TerminalAborted,
			Stats:     agent.RunStats{},
		}
	}()
	return events
}
