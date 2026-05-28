package tool

import (
	"context"
	"errors"
	"fmt"
)

const defaultRuntimeRegistryLimit = 20

// RuntimeRegistry stores tools discovered while a session is active.
type RuntimeRegistry struct {
	tools          map[string]Tool
	definitions    map[string]Definition
	order          []string
	lastAccessTick map[string]uint64
	accessTick     uint64
	limit          int
}

// NewRuntimeRegistry creates a bounded registry for runtime-discovered tools.
func NewRuntimeRegistry(limit int) *RuntimeRegistry {
	if limit <= 0 {
		limit = defaultRuntimeRegistryLimit
	}

	return &RuntimeRegistry{
		tools:          make(map[string]Tool),
		definitions:    make(map[string]Definition),
		order:          []string{},
		lastAccessTick: make(map[string]uint64),
		limit:          limit,
	}
}

// Register adds one runtime-discovered tool and evicts least-recently-used
// tools when the registry exceeds its configured limit.
func (r *RuntimeRegistry) Register(candidate Tool) error {
	if r == nil {
		return errors.New("register runtime tool: nil registry")
	}
	if candidate == nil {
		return errors.New("register runtime tool: nil tool")
	}

	definition := candidate.Definition()
	if definition.Name == "" {
		return errors.New("register runtime tool: empty tool name")
	}
	if _, exists := r.tools[definition.Name]; exists {
		return fmt.Errorf("register runtime tool %q: duplicate tool name", definition.Name)
	}

	r.tools[definition.Name] = candidate
	r.definitions[definition.Name] = definition
	r.order = append(r.order, definition.Name)
	r.markAccessed(definition.Name)
	r.evictOverflow()
	return nil
}

// Has reports whether a runtime tool name is registered.
func (r *RuntimeRegistry) Has(name string) bool {
	if r == nil {
		return false
	}

	_, exists := r.tools[name]
	return exists
}

// Definition returns one runtime tool definition by name.
func (r *RuntimeRegistry) Definition(name string) (Definition, bool) {
	if r == nil {
		return Definition{}, false
	}

	definition, exists := r.definitions[name]
	return definition, exists
}

// Definitions returns runtime tool definitions in registration order.
func (r *RuntimeRegistry) Definitions() []Definition {
	if r == nil {
		return nil
	}

	definitions := make([]Definition, 0, len(r.order))
	for _, name := range r.order {
		definitions = append(definitions, r.definitions[name])
	}
	return definitions
}

// Unregister removes one runtime tool and its definition.
func (r *RuntimeRegistry) Unregister(name string) bool {
	if r == nil {
		return false
	}
	if _, exists := r.tools[name]; !exists {
		return false
	}

	delete(r.tools, name)
	delete(r.definitions, name)
	delete(r.lastAccessTick, name)
	for index, registeredName := range r.order {
		if registeredName == name {
			r.order = append(r.order[:index], r.order[index+1:]...)
			break
		}
	}
	return true
}

// MarkAccessed records recent runtime tool usage for LRU eviction.
func (r *RuntimeRegistry) MarkAccessed(name string) {
	if r == nil {
		return
	}
	if _, exists := r.tools[name]; !exists {
		return
	}

	r.markAccessed(name)
}

// Execute finds the requested runtime tool by name and runs it.
func (r *RuntimeRegistry) Execute(ctx context.Context, call Call) (Result, error) {
	if r == nil {
		return Result{}, fmt.Errorf("execute runtime tool %q: %w", call.Name, ErrNotFound)
	}

	selected, exists := r.tools[call.Name]
	if !exists {
		return Result{}, fmt.Errorf("execute runtime tool %q: %w", call.Name, ErrNotFound)
	}

	result, err := selected.Execute(ctx, call)
	if err != nil {
		return Result{}, fmt.Errorf("execute runtime tool %q: %w", call.Name, err)
	}
	r.MarkAccessed(call.Name)
	return result, nil
}

func (r *RuntimeRegistry) evictOverflow() {
	for len(r.order) > r.limit {
		r.Unregister(r.leastRecentlyUsedName())
	}
}

func (r *RuntimeRegistry) leastRecentlyUsedName() string {
	if len(r.order) == 0 {
		return ""
	}

	oldestName := r.order[0]
	oldestTick := r.lastAccessTick[oldestName]
	for _, name := range r.order[1:] {
		accessTick := r.lastAccessTick[name]
		if accessTick < oldestTick {
			oldestName = name
			oldestTick = accessTick
		}
	}
	return oldestName
}

func (r *RuntimeRegistry) markAccessed(name string) {
	r.accessTick++
	r.lastAccessTick[name] = r.accessTick
}
