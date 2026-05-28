package app

import (
	"context"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

// Runtime owns the reusable agent core shared by TUI and eval frontends.
type Runtime struct {
	State   *State
	Agent   agent.Agent
	Session *session.State
}

// NewRuntime wires provider, memory, tools, prompt, and session state.
func NewRuntime(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	provider, err := newProvider(cfg, logger)
	if err != nil {
		return nil, err
	}

	memoryStore := newMemoryStore(cfg)

	mcpManager, closers, err := newMCPManager(ctx, cfg, logger)
	if err != nil {
		logger.DebugContext(ctx, "MCP manager initialization failed", "error", err)
		return nil, err
	}

	state := &State{
		Config:      cfg,
		Logger:      logger,
		Provider:    provider,
		MemoryStore: memoryStore,
		MCPManager:  mcpManager,
		Closers:     closers,
	}

	sessionState, err := newSession(ctx, state)
	if err != nil {
		logger.DebugContext(ctx, "Session initialization failed", "error", err)
		closeAll(ctx, logger, closers)
		return nil, err
	}
	logger.DebugContext(ctx, "Session toolset initialized", "tool_count", len(sessionState.ToolDefinitions()))

	runtimeAgent, err := newAgent(cfg, provider, logger)
	if err != nil {
		closeAll(ctx, logger, closers)
		return nil, err
	}

	return &Runtime{
		State:   state,
		Agent:   runtimeAgent,
		Session: sessionState,
	}, nil
}

// Close releases runtime-owned resources such as MCP sessions.
func (r *Runtime) Close(ctx context.Context) {
	if r == nil || r.State == nil {
		return
	}
	closeAll(ctx, r.State.Logger, r.State.Closers)
}
