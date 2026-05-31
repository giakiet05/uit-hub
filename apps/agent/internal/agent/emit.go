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
