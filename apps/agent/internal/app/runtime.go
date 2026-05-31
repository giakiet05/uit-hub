package app

import (
	"context"
	"log/slog"


	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eventbus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eventhandler"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

// Runtime owns the reusable agent core shared by TUI and eval frontends.
type Runtime struct {
	State              *State
	Session            *session.Session
	Bus                *eventbus.EventBus
	UsageEventHandler  *eventhandler.UsageEventHandler
	LoggerEventHandler *eventhandler.LoggerEventHandler
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

	sessionState, err := newSessionState(ctx, state)
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

	bus := eventbus.New(128, logger)
	
	usageSub := eventhandler.NewUsageEventHandler(ctx, sessionState, bus, logger)
	usageSub.Start()
	
	loggerSub := eventhandler.NewLoggerEventHandler(ctx, bus, logger)
	loggerSub.Start()

	logger.DebugContext(ctx, "Startup configuration",
		"provider", cfg.Provider,
		"agent_type", cfg.Agent.Type,
		"agent_max_rounds", cfg.Agent.MaxRounds,
		"agent_concurrent_tools", cfg.Agent.ConcurrentTools,
		"agent_tool_timeout", cfg.Agent.ToolTimeout.String(),
		"mcp_servers", len(cfg.MCP.Servers),
		"session_tools", len(sessionState.ToolDefinitions()),
	)

	return &Runtime{
		State:              state,
		Session:            session.NewSession(sessionState, runtimeAgent, bus),
		Bus:                bus,
		UsageEventHandler:  usageSub,
		LoggerEventHandler: loggerSub,
	}, nil
}

func (r *Runtime) Close(ctx context.Context) {
	if r == nil {
		return
	}
	if r.UsageEventHandler != nil {
		r.UsageEventHandler.Stop()
	}
	if r.LoggerEventHandler != nil {
		r.LoggerEventHandler.Stop()
	}
	if r.Bus != nil {
		r.Bus.Close()
	}
	if r.State != nil {
		closeAll(ctx, r.State.Logger, r.State.Closers)
	}
}
