package tool

import (
	"context"
	"fmt"
)

// ToolSet exposes base tools and runtime-discovered tools as one executable
// collection to agents.
type ToolSet struct {
	base    *BaseRegistry
	runtime *RuntimeRegistry
}

// NewToolSet creates a combined tool collection.
func NewToolSet(base *BaseRegistry, runtime *RuntimeRegistry) *ToolSet {
	return &ToolSet{
		base:    base,
		runtime: runtime,
	}
}

// Definitions returns all tool definitions visible to the model.
func (s *ToolSet) Definitions() []Definition {
	if s == nil {
		return nil
	}

	baseDefinitions := s.base.Definitions()
	runtimeDefinitions := s.runtime.Definitions()
	definitions := make([]Definition, 0, len(baseDefinitions)+len(runtimeDefinitions))
	definitions = append(definitions, baseDefinitions...)
	definitions = append(definitions, runtimeDefinitions...)
	return definitions
}

// Definition returns one visible tool definition by name.
func (s *ToolSet) Definition(name string) (Definition, bool) {
	if s == nil {
		return Definition{}, false
	}
	if definition, exists := s.base.Definition(name); exists {
		return definition, true
	}
	return s.runtime.Definition(name)
}

// Tool returns one visible tool by name, preferring base tools over runtime
// tools when names collide.
func (s *ToolSet) Tool(name string) (Tool, bool) {
	if s == nil {
		return nil, false
	}
	if selected, exists := s.base.Tool(name); exists {
		return selected, true
	}
	return s.runtime.Tool(name)
}

// Execute runs a tool from the base registry first, then the runtime registry.
func (s *ToolSet) Execute(ctx context.Context, call Call) (Result, error) {
	if s == nil {
		return Result{}, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound)
	}

	if s.base.Has(call.Name) {
		return s.base.Execute(ctx, call)
	}
	if s.runtime.Has(call.Name) {
		return s.runtime.Execute(ctx, call)
	}
	return Result{}, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound)
}

// MarkUsed records recent runtime tool usage when the tool is runtime-owned.
func (s *ToolSet) MarkUsed(name string) {
	if s == nil {
		return
	}
	s.runtime.MarkAccessed(name)
}

// RuntimeToolNames returns the runtime-discovered tool names currently visible
// to the model.
func (s *ToolSet) RuntimeToolNames() []string {
	if s == nil {
		return nil
	}
	return s.runtime.Names()
}
