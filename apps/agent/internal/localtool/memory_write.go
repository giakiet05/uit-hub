package localtool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// MemoryWrite creates or updates long-term memory.
type MemoryWrite struct {
	tool.BaseTool
	store memory.Store
}

// NewMemoryWrite creates a memory writing tool.
func NewMemoryWrite(store memory.Store) *MemoryWrite {
	return &MemoryWrite{
		BaseTool: tool.NewBaseTool(
			tool.Definition{
				Name:        "memory_write",
				Description: "Create or update durable long-term memory only for stable user preferences, durable feedback, project context, or references. Do not save transient task data, student records, grades, or backend-source data.",
				InputSchema: tool.ObjectSchema(
					map[string]any{
						"id":          tool.StringProperty("Existing memory ID to update, or empty to create a new memory."),
						"type":        tool.StringEnumProperty("Memory type.", "user", "feedback", "project", "reference"),
						"name":        tool.StringProperty("Short human-readable memory name."),
						"description": tool.StringProperty("One-line relevance description for future recall."),
						"content":     tool.StringProperty("Full memory content in Markdown."),
					},
					"type",
					"name",
					"description",
					"content",
				),
			},
			tool.NewWriteMetadata(false, tool.DefaultMaxResultChars),
		),
		store: store,
	}
}

// Execute creates or updates a memory and returns the saved memory metadata as
// JSON.
func (t *MemoryWrite) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}
	if t.store == nil {
		return tool.Result{}, memory.ErrDisabled
	}

	id, err := optionalStringArg(call.Arguments, "id")
	if err != nil {
		return tool.Result{}, err
	}
	memoryType, err := stringArg(call.Arguments, "type")
	if err != nil {
		return tool.Result{}, err
	}
	name, err := stringArg(call.Arguments, "name")
	if err != nil {
		return tool.Result{}, err
	}
	description, err := stringArg(call.Arguments, "description")
	if err != nil {
		return tool.Result{}, err
	}
	content, err := stringArg(call.Arguments, "content")
	if err != nil {
		return tool.Result{}, err
	}

	saved, err := t.store.Write(ctx, memory.WriteInput{
		ID:          id,
		Type:        memory.Type(memoryType),
		Name:        name,
		Description: description,
		Content:     content,
	})
	if err != nil {
		return tool.Result{}, err
	}
	output, err := json.Marshal(map[string]any{
		"id":          saved.ID,
		"type":        saved.Type,
		"name":        saved.Name,
		"description": saved.Description,
		"updated_at":  saved.UpdatedAt,
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal saved memory: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}
