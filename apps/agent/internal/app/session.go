package app

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

// newSession creates one chat session with stable memory, prompt, and tool
// snapshots.
func newSession(ctx context.Context, state *State) (*session.State, error) {
	memoryContext, err := memory.NewLoader(state.MemoryStore).Load(ctx)
	if err != nil {
		state.Logger.DebugContext(ctx, "Memory context load failed", "error", err)
		return nil, err
	}

	tools, err := newSessionToolSet(state.Config, state.MemoryStore, state.MCPManager)
	if err != nil {
		return nil, err
	}

	promptSnapshot := newSessionPrompt(memoryContext, state.MCPManager.Catalog())
	return session.NewState(
		session.WithPromptSnapshot(promptSnapshot),
		session.WithTools(tools),
	), nil
}
