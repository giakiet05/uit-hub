package app

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
)

// newMemoryStore creates the current memory backend. FileStore is temporary CLI
// storage.
func newMemoryStore(cfg config.Config) memory.Store {
	return memory.NewFileStore(cfg.Memory.Path)
}
