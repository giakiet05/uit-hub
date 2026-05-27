package app

import (
	"context"
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

// Runtime owns the reusable agent core shared by TUI and eval frontends.
type Runtime struct {
	Agent   agent.Agent
	Session *runtime.Session
	closers []io.Closer
	logger  *slog.Logger
}

// NewRuntime wires provider, memory, tools, prompt, and session state.
func NewRuntime(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	provider, err := newProvider(cfg, logger)
	if err != nil {
		return nil, err
	}

	session := runtime.NewSession()
	memoryStore := newMemoryStore(cfg)
	memoryContext, err := memory.NewLoader(memoryStore).Load(ctx)
	if err != nil {
		logger.DebugContext(ctx, "Memory context load failed", "error", err)
		return nil, err
	}

	tools, mcpManager, closers, err := newToolSet(ctx, cfg, logger, memoryStore)
	if err != nil {
		logger.DebugContext(ctx, "Tool registry initialization failed", "error", err)
		return nil, err
	}
	logger.DebugContext(ctx, "Tool registry initialized", "tool_count", len(tools.Definitions()))

	promptBuilder := prompt.NewBuilder(newSessionPrompt(memoryContext, mcpManager.Catalog()))
	runtimeAgent, err := newAgent(cfg, provider, promptBuilder, tools, logger)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Agent:   runtimeAgent,
		Session: session,
		closers: closers,
		logger:  logger,
	}, nil
}

// Close releases runtime-owned resources such as MCP sessions.
func (r *Runtime) Close(ctx context.Context) {
	if r == nil {
		return
	}
	closeAll(ctx, r.logger, r.closers)
}
