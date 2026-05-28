// Package agent defines the runtime contract and concrete agent loops.
package agent

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

// Agent runs one user prompt against a session and streams loop events until a
// terminal event is emitted and the channel closes.
type Agent interface {
	Run(ctx context.Context, session *session.State, input string) <-chan Event
}
