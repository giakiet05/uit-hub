package localtool

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// MCPToolLoader loads a full MCP tool wrapper from the compact catalog.
type MCPToolLoader interface {
	LoadTool(name string) (tool.Tool, bool)
}

// RuntimeToolRegistry registers runtime-discovered tools.
type RuntimeToolRegistry interface {
	Has(name string) bool
	Register(candidate tool.Tool) error
}

// LoadMCPTool loads one MCP catalog tool into the runtime registry.
type LoadMCPTool struct {
	loader   MCPToolLoader
	registry RuntimeToolRegistry
}

// NewLoadMCPTool creates a tool that makes deferred MCP tools available to the
// model in the next round.
func NewLoadMCPTool(loader MCPToolLoader, registry RuntimeToolRegistry) *LoadMCPTool {
	return &LoadMCPTool{
		loader:   loader,
		registry: registry,
	}
}

// Definition describes the load_mcp_tool schema.
func (t *LoadMCPTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "load_mcp_tool",
		Description: "Load one MCP catalog tool by exact namespaced name so it can be used in the next round.",
		InputSchema: tool.ObjectSchema(
			map[string]any{
				"name": tool.StringProperty("Exact namespaced MCP tool name from the MCP tool catalog."),
			},
			"name",
		),
	}
}

// Execute loads one MCP tool into the runtime registry.
func (t *LoadMCPTool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}
	if t.loader == nil {
		return tool.Result{}, fmt.Errorf("mcp tool loader is not configured")
	}
	if t.registry == nil {
		return tool.Result{}, fmt.Errorf("runtime tool registry is not configured")
	}

	name, err := stringArg(call.Arguments, "name")
	if err != nil {
		return tool.Result{}, err
	}
	if t.registry.Has(name) {
		return tool.Result{
			CallID:  call.ID,
			Name:    call.Name,
			Content: fmt.Sprintf("MCP tool %s is already loaded.", name),
		}, nil
	}

	loadedTool, exists := t.loader.LoadTool(name)
	if !exists {
		return tool.Result{}, fmt.Errorf("unknown MCP tool %q", name)
	}
	if err := t.registry.Register(loadedTool); err != nil {
		return tool.Result{}, err
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: fmt.Sprintf("MCP tool %s loaded. Use it in the next round.", name),
	}, nil
}
