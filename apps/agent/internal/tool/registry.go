package tool

import (
	"context"
	"errors"
	"fmt"
)

// ErrNotFound marks a request for a tool name that is not registered.
var ErrNotFound = errors.New("tool not found")

// Registry stores tools by name and preserves registration order for schema
// generation.
type Registry struct {
	tools       map[string]Tool
	definitions map[string]Definition
	order       []string
}

// NewRegistry creates a registry and registers the provided tools.
func NewRegistry(tools ...Tool) (*Registry, error) {
	registry := &Registry{
		tools:       make(map[string]Tool, len(tools)),
		definitions: make(map[string]Definition, len(tools)),
		order:       make([]string, 0, len(tools)),
	}

	for _, candidate := range tools {
		if err := registry.Register(candidate); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

// Register adds one tool to the registry after validating its definition.
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
	r.definitions[definition.Name] = definition
	r.order = append(r.order, definition.Name)
	return nil
}

// Definitions returns tool definitions in registration order.
func (r *Registry) Definitions() []Definition {
	if r == nil {
		return nil
	}

	definitions := make([]Definition, 0, len(r.order))
	for _, name := range r.order {
		definitions = append(definitions, r.definitions[name])
	}
	return definitions
}

// Execute finds the requested tool by name and runs it.
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
