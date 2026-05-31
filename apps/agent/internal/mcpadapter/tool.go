package mcpadapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool adapts one MCP server tool to the agent's native tool interface.
type Tool struct {
	tool.BaseTool
	serverName string
	session    Session
	mcpTool    *mcp.Tool
}

// NewTool creates a native wrapper for one MCP tool.
func NewTool(serverName string, session Session, mcpTool *mcp.Tool) *Tool {
	concurrencySafe := false
	requireApproval := false

	if mcpTool.InputSchema != nil {
		data, err := json.Marshal(mcpTool.InputSchema)
		if err == nil {
			var m map[string]any
			if json.Unmarshal(data, &m) == nil {
				if val, ok := m["_concurrencySafe"].(bool); ok && val {
					concurrencySafe = true
				}
				delete(m, "_concurrencySafe")
				
				if val, ok := m["_requireApproval"].(bool); ok && val {
					requireApproval = true
				}
				delete(m, "_requireApproval")
				
				mcpTool.InputSchema = m
			}
		}
	}

	return &Tool{
		BaseTool: tool.NewBaseTool(
			tool.Definition{
				Name:        NamespacedName(serverName, mcpTool.Name),
				Description: mcpTool.Description,
				InputSchema: InputSchemaFromMCP(mcpTool.InputSchema),
			},
			tool.Metadata{
				ReadOnly:        false,
				Destructive:     false,
				ConcurrencySafe: concurrencySafe,
				RequireApproval: requireApproval,
				MaxResultChars:  tool.DefaultMaxResultChars,
			},
		),
		serverName: serverName,
		session:    session,
		mcpTool:    mcpTool,
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
