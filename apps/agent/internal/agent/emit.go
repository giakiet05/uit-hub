package agent

import "context"

// Emit sends one agent event unless the context is canceled.
func Emit(ctx context.Context, events chan<- Event, event Event) bool {
	select {
	case <-ctx.Done():
		return false
	case events <- event:
		return true
	}
}

// EmitTerminal sends a terminal or cleanup event even when the run context was
// canceled. Use it only for events needed to close and persist a run cleanly.
func EmitTerminal(events chan<- Event, event Event) bool {
	events <- event
	return true
}
