// Package agent defines the runtime contract and concrete agent loops.
package agent

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

// Agent runs one user prompt against a session and returns the assistant message
// produced by the selected agent loop.
type Agent interface {
	Run(ctx context.Context, session *runtime.Session, input string) (conversation.Message, error)
}
