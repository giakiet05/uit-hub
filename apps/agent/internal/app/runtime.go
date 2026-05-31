package app

import (
	"context"
	"log/slog"


	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/handler"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

// Runtime owns the reusable agent core shared by TUI and eval frontends.
type Runtime struct {
	State              *State
	Session            *session.Session
	Bus                *bus.EventBus
	UsageEventHandler   *handler.UsageEventHandler
	LoggerEventHandler  *handler.LoggerEventHandler
	StorageEventHandler *handler.StorageEventHandler
}

// NewRuntime wires provider, memory, tools, prompt, and session state.
func NewRuntime(ctx context.Context, cfg config.Config, logger *slog.Logger, resume bool, resumeID string) (*Runtime, error) {
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

	db, err := storage.InitDB("agent.db")
	if err != nil {
		logger.ErrorContext(ctx, "Failed to initialize SQLite database", "error", err)
	}

	var sessionOpts []session.Option
	if db != nil && resume {
		var record *storage.SessionRecord
		var dbErr error
		if resumeID != "" {
			record, dbErr = db.GetSession(resumeID)
		} else {
			record, dbErr = db.GetLastSession()
		}
		if dbErr == nil && record != nil {
			messages, usage, decErr := record.Decode()
			if decErr == nil {
				sessionOpts = append(sessionOpts, session.WithHistory(record.ID, record.StartedAt, messages, usage))
			} else {
				logger.ErrorContext(ctx, "Failed to decode session history", "error", decErr)
			}
		} else if dbErr != nil {
			logger.ErrorContext(ctx, "Failed to load session", "error", dbErr)
		}
	}

	sessionState, err := newSessionState(ctx, state, sessionOpts...)
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

	bus := bus.New(128, logger)
	
	usageSub := handler.NewUsageEventHandler(ctx, sessionState, bus, logger)
	usageSub.Start()
	
	loggerSub := handler.NewLoggerEventHandler(ctx, bus, logger)
	loggerSub.Start()



	var storageSub *handler.StorageEventHandler
	if db != nil {
		storageSub = handler.NewStorageEventHandler(ctx, db, sessionState, bus, logger)
		storageSub.Start()
	}

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
		Bus:                 bus,
		UsageEventHandler:   usageSub,
		LoggerEventHandler:  loggerSub,
		StorageEventHandler: storageSub,
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
	if r.StorageEventHandler != nil {
		r.StorageEventHandler.Stop()
	}
	if r.Bus != nil {
		r.Bus.Close()
	}
	if r.State != nil {
		closeAll(ctx, r.State.Logger, r.State.Closers)
	}
}
