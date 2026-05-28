package app

import (
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
)

// State owns process-lifetime dependencies shared by sessions and frontends.
type State struct {
	Config      config.Config
	Logger      *slog.Logger
	Provider    llm.Provider
	MemoryStore memory.Store
	MCPManager  *mcpadapter.Manager
	Closers     []io.Closer
}
