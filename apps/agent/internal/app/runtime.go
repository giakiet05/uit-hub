package app

import (
	"context"
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

// Runtime owns the reusable agent core shared by TUI and eval frontends.
type Runtime struct {
	State              *State
	Session            *session.Session
	Bus                *bus.EventBus
	EventHandlerCloser io.Closer
}

// NewRuntime wires memory, tools, prompt, and session state.
func NewRuntime(ctx context.Context, cfg config.Config, logger *slog.Logger, model llm.Model, resume bool, resumeID string) (*Runtime, error) {

	memoryStore := newMemoryStore(cfg)

	mcpManager, closers, err := newMCPManager(ctx, cfg, logger)
	if err != nil {
		logger.DebugContext(ctx, "MCP manager initialization failed", "error", err)
		return nil, err
	}

	state := &State{
		Config:      cfg,
		Logger:      logger,
		MemoryStore: memoryStore,
		MCPManager:  mcpManager,
		Closers:     closers,
	}

	db, err := storage.InitDB(cfg.Storage.DBPath)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to initialize SQLite database", "error", err)
	}

	var sessionOpts []session.Option
	var runtimeToolNames []string
	if db != nil && resume {
		var record *storage.SessionRecord
		var dbErr error
		if resumeID != "" {
			record, dbErr = db.GetSession(resumeID)
		} else {
			record, dbErr = db.GetLastSession()
		}
		if dbErr == nil && record != nil {
			messages, usage, restoredRuntimeToolNames, decErr := record.Decode()
			if decErr == nil {
				sessionOpts = append(sessionOpts, session.WithHistory(record.ID, record.StartedAt, messages, usage))
				runtimeToolNames = restoredRuntimeToolNames
			} else {
				logger.ErrorContext(ctx, "Failed to decode session history", "error", decErr)
			}
		} else if dbErr != nil {
			logger.ErrorContext(ctx, "Failed to load session", "error", dbErr)
		}
	}

	sessionState, err := newSessionState(ctx, state, model, runtimeToolNames, sessionOpts...)
	if err != nil {
		logger.DebugContext(ctx, "Session initialization failed", "error", err)
		closeAll(ctx, logger, closers)
		return nil, err
	}
	logger.DebugContext(ctx, "Session toolset initialized", "tool_count", len(sessionState.ToolDefinitions()))

	bus := bus.New(128, logger)
	handlerCloser := StartEventHandlers(ctx, bus, logger)

	logger.DebugContext(ctx, "Startup configuration",
		"agent_type", cfg.Agent.Type,
		"agent_max_rounds", cfg.Agent.MaxRounds,
		"agent_concurrent_tools", cfg.Agent.ConcurrentTools,
		"agent_tool_timeout", cfg.Agent.ToolTimeout.String(),
		"mcp_servers", len(cfg.MCP.Servers),
		"session_tools", len(sessionState.ToolDefinitions()),
	)

	return &Runtime{
		State:              state,
		Session:            session.NewSession(sessionState, nil, bus, session.WithStore(db), session.WithLogger(logger), session.WithShutdownContext(ctx)),
		Bus:                bus,
		EventHandlerCloser: handlerCloser,
	}, nil
}

func (r *Runtime) Close(ctx context.Context) {
	if r == nil {
		return
	}
	if r.EventHandlerCloser != nil {
		r.EventHandlerCloser.Close()
	}
	if r.Bus != nil {
		r.Bus.Close()
	}
	if r.State != nil {
		closeAll(ctx, r.State.Logger, r.State.Closers)
	}
}
