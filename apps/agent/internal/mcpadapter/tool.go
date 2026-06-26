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

				// Strip "token" from schema to hide it from LLM
				if props, ok := m["properties"].(map[string]any); ok {
					delete(props, "token")
				}
				if reqs, ok := m["required"].([]any); ok {
					newReqs := make([]any, 0, len(reqs))
					for _, req := range reqs {
						if r, ok := req.(string); ok && r == "token" {
							continue
						}
						newReqs = append(newReqs, req)
					}
					if len(newReqs) == 0 {
						delete(m, "required")
					} else {
						m["required"] = newReqs
					}
				}

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
				ReadOnly:          false,
				Destructive:       false,
				ConcurrencySafe:   concurrencySafe,
				RequireApproval:   requireApproval,
				FinishOnInterrupt: requireApproval,
			},
		),
		serverName: serverName,
		session:    session,
		mcpTool:    mcpTool,
	}
}

// Execute calls the underlying MCP tool and flattens the result for the model.
func (t *Tool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	// Inject token from context if present
	if token, ok := ctx.Value(TokenKey).(string); ok && token != "" {
		if call.Arguments == nil {
			call.Arguments = make(map[string]any)
		}
		call.Arguments["token"] = token
	}

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
