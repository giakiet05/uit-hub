package tool

import (
	"context"
	"errors"
	"fmt"
)

// ErrNotFound marks a request for a tool name that is not registered.
var ErrNotFound = errors.New("tool not found")

// BaseRegistry stores preloaded tools by name and preserves registration order
// for schema generation.
type BaseRegistry struct {
	tools map[string]Tool
	// definitions snapshots tool metadata at registration time so later LLM
	// calls see stable schemas even if a tool builds definitions dynamically.
	definitions map[string]Definition
	order       []string
}

// NewBaseRegistry creates a base registry and registers the provided tools.
func NewBaseRegistry(tools ...Tool) (*BaseRegistry, error) {
	registry := &BaseRegistry{
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
func (r *BaseRegistry) Register(candidate Tool) error {
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

// Has reports whether a tool name is registered.
func (r *BaseRegistry) Has(name string) bool {
	if r == nil {
		return false
	}

	_, exists := r.tools[name]
	return exists
}

// Definition returns one registered tool definition by name.
func (r *BaseRegistry) Definition(name string) (Definition, bool) {
	if r == nil {
		return Definition{}, false
	}

	definition, exists := r.definitions[name]
	return definition, exists
}

// Definitions returns tool definitions in registration order.
func (r *BaseRegistry) Definitions() []Definition {
	if r == nil {
		return nil
	}

	definitions := make([]Definition, 0, len(r.order))
	for _, name := range r.order {
		definitions = append(definitions, r.definitions[name])
	}
	return definitions
}

// Unregister removes one tool and its definition from the registry.
func (r *BaseRegistry) Unregister(name string) bool {
	if r == nil {
		return false
	}
	if _, exists := r.tools[name]; !exists {
		return false
	}

	delete(r.tools, name)
	delete(r.definitions, name)
	for index, registeredName := range r.order {
		if registeredName == name {
			r.order = append(r.order[:index], r.order[index+1:]...)
			break
		}
	}
	return true
}

// Execute finds the requested tool by name and runs it.
func (r *BaseRegistry) Execute(ctx context.Context, call Call) (Result, error) {
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
