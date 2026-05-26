package mcpadapter

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool adapts one MCP server tool to the agent's native tool interface.
type Tool struct {
	serverName string
	session    Session
	mcpTool    *mcp.Tool
}

// NewTool creates a native wrapper for one MCP tool.
func NewTool(serverName string, session Session, mcpTool *mcp.Tool) *Tool {
	return &Tool{
		serverName: serverName,
		session:    session,
		mcpTool:    mcpTool,
	}
}

// Definition returns the native tool definition passed to LLM providers.
func (t *Tool) Definition() tool.Definition {
	return tool.Definition{
		Name:        NamespacedName(t.serverName, t.mcpTool.Name),
		Description: t.mcpTool.Description,
		InputSchema: InputSchemaFromMCP(t.mcpTool.InputSchema),
	}
}

// Execute calls the underlying MCP tool and flattens the result for the model.
func (t *Tool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	result, err := t.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      t.mcpTool.Name,
		Arguments: call.Arguments,
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("call mcp tool %q: %w", t.mcpTool.Name, err)
	}

	content, err := FlattenToolResult(result)
	if err != nil {
		return tool.Result{}, err
	}
	if result.IsError {
		content = "mcp tool error: " + content
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: content,
	}, nil
}
