package localtool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// MemoryList lists saved memory index entries.
type MemoryList struct {
	store memory.Store
}

// NewMemoryList creates a memory listing tool.
func NewMemoryList(store memory.Store) *MemoryList {
	return &MemoryList{store: store}
}

// Definition describes the memory_list tool schema.
func (t *MemoryList) Definition() tool.Definition {
	return tool.Definition{
		Name:        "memory_list",
		Description: "List saved long-term memories by id, type, name, and description.",
		InputSchema: tool.EmptyInputSchema(),
	}
}

// Execute returns memory index entries as JSON.
func (t *MemoryList) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}
	if t.store == nil {
		return tool.Result{}, memory.ErrDisabled
	}

	entries, err := t.store.List(ctx)
	if err != nil {
		return tool.Result{}, err
	}
	output, err := json.Marshal(entries)
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal memory list: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}
