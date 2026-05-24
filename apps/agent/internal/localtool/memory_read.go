package localtool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// MemoryRead reads one saved memory by ID.
type MemoryRead struct {
	store memory.Store
}

// NewMemoryRead creates a memory reading tool.
func NewMemoryRead(store memory.Store) *MemoryRead {
	return &MemoryRead{store: store}
}

// Definition describes the memory_read tool schema.
func (t *MemoryRead) Definition() tool.Definition {
	return tool.Definition{
		Name:        "memory_read",
		Description: "Read one saved long-term memory by id when the memory index indicates it may be relevant.",
		InputSchema: tool.ObjectSchema(
			map[string]any{
				"id": tool.StringProperty("Memory ID to read."),
			},
			"id",
		),
	}
}

// Execute returns the full memory as JSON.
func (t *MemoryRead) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}
	if t.store == nil {
		return tool.Result{}, memory.ErrDisabled
	}

	id, err := stringArg(call.Arguments, "id")
	if err != nil {
		return tool.Result{}, err
	}
	selected, err := t.store.Read(ctx, id)
	if err != nil {
		return tool.Result{}, err
	}
	output, err := json.Marshal(selected)
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal memory: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}
