package tool

import (
	"context"
	"errors"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

var ErrNotFound = errors.New("tool not found")

type Registry struct {
	tools map[string]Tool
	order []string
}

func NewRegistry(tools ...Tool) (*Registry, error) {
	registry := &Registry{
		tools: make(map[string]Tool, len(tools)),
		order: make([]string, 0, len(tools)),
	}

	for _, candidate := range tools {
		if err := registry.Register(candidate); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

func (r *Registry) Register(candidate Tool) error {
	if candidate == nil {
		return errors.New("register tool: nil tool")
	}

	definition := candidate.Definition()
	if definition.Name == "" {
		return errors.New("register tool: empty tool name")
	}
	if _, exists := r.tools[definition.Name]; exists {
		return fmt.Errorf("register tool %q: duplicate tool name", definition.Name)
	}

	r.tools[definition.Name] = candidate
	r.order = append(r.order, definition.Name)
	return nil
}

func (r *Registry) Definitions() []llm.ToolDefinition {
	if r == nil {
		return nil
	}

	definitions := make([]llm.ToolDefinition, 0, len(r.order))
	for _, name := range r.order {
		definitions = append(definitions, r.tools[name].Definition())
	}
	return definitions
}

func (r *Registry) Execute(ctx context.Context, call Call) (Result, error) {
	if r == nil {
		return Result{}, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound)
	}

	selected, exists := r.tools[call.Name]
	if !exists {
		return Result{}, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound)
	}

	result, err := selected.Execute(ctx, call)
	if err != nil {
		return Result{}, fmt.Errorf("execute tool %q: %w", call.Name, err)
	}
	return result, nil
}
