package agent

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

type Agent interface {
	Run(ctx context.Context, session *runtime.Session, input string) (conversation.Message, error)
}
