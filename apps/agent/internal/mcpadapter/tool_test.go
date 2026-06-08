package mcpadapter

import (
	"context"
	"errors"
	"strings"
	"testing"

	agenttool "github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestClientToolsWrapsMCPTools(t *testing.T) {
	session := &fakeSession{
		tools: []*mcp.Tool{
			{
				Name:        "get_student",
				Description: "Get student profile.",
				InputSchema: map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]any{
						"id": map[string]any{"type": "string"},
					},
				},
			},
		},
	}
	client := NewClient("uit", session)

	tools, err := client.Tools(context.Background())
	if err != nil {
		t.Fatalf("Tools() error = %v", err)
	}
	if got, want := len(tools), 1; got != want {
		t.Fatalf("tool count = %d, want %d", got, want)
	}

	definition := tools[0].Definition()
	if got, want := definition.Name, "uit__get_student"; got != want {
		t.Fatalf("definition name = %q, want %q", got, want)
	}
	if got, want := definition.Description, "Get student profile."; got != want {
		t.Fatalf("definition description = %q, want %q", got, want)
	}
	inputSchema := definition.InputSchema
	properties, ok := inputSchema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema properties = %T, want map[string]any", inputSchema["properties"])
	}
	idProperty, ok := properties["id"].(map[string]any)
	if !ok {
		t.Fatalf("id property = %T, want map[string]any", properties["id"])
	}
	if got, want := idProperty["type"], "string"; got != want {
		t.Fatalf("schema property type = %v, want %v", got, want)
	}
	if tools[0].Metadata().FinishOnInterrupt {
		t.Fatal("read-only MCP tool should not finish on interrupt")
	}
}

func TestClientToolsMarksApprovalToolsAsFinishOnInterrupt(t *testing.T) {
	session := &fakeSession{
		tools: []*mcp.Tool{
			{
				Name:        "send_email",
				Description: "Send an email.",
				InputSchema: map[string]any{
					"type":             "object",
					"_requireApproval": true,
				},
			},
		},
	}
	client := NewClient("gmail", session)

	tools, err := client.Tools(context.Background())
	if err != nil {
		t.Fatalf("Tools() error = %v", err)
	}

	metadata := tools[0].Metadata()
	if !metadata.RequireApproval {
		t.Fatal("RequireApproval = false, want true")
	}
	if !metadata.FinishOnInterrupt {
		t.Fatal("FinishOnInterrupt = false, want true")
	}
}

func TestToolExecuteCallsMCPTool(t *testing.T) {
	session := &fakeSession{
		result: &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "student profile"},
			},
			StructuredContent: map[string]any{
				"ok": true,
			},
		},
	}
	adapter := NewTool("uit", session, &mcp.Tool{Name: "get_student"})

	result, err := adapter.Execute(context.Background(), agenttool.Call{
		ID:   "call-1",
		Name: "uit__get_student",
		Arguments: map[string]any{
			"id": "22520001",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := session.calledName, "get_student"; got != want {
		t.Fatalf("called tool = %q, want %q", got, want)
	}
	if got, want := result.CallID, "call-1"; got != want {
		t.Fatalf("result call id = %q, want %q", got, want)
	}
	if got, want := result.Name, "uit__get_student"; got != want {
		t.Fatalf("result name = %q, want %q", got, want)
	}
	if !strings.Contains(result.Content, "student profile") {
		t.Fatalf("result content missing text: %q", result.Content)
	}
	if !strings.Contains(result.Content, `"ok":true`) {
		t.Fatalf("result content missing structured content: %q", result.Content)
	}
}

func TestToolExecutePrefixesMCPToolError(t *testing.T) {
	session := &fakeSession{
		result: &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "not found"},
			},
		},
	}
	adapter := NewTool("uit", session, &mcp.Tool{Name: "get_student"})

	result, err := adapter.Execute(context.Background(), agenttool.Call{
		ID:   "call-1",
		Name: "uit__get_student",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := result.Content, "mcp tool error: not found"; got != want {
		t.Fatalf("result content = %q, want %q", got, want)
	}
}

func TestSplitNamespacedName(t *testing.T) {
	serverName, toolName, err := SplitNamespacedName("uit__get_student")
	if err != nil {
		t.Fatalf("SplitNamespacedName() error = %v", err)
	}
	if serverName != "uit" || toolName != "get_student" {
		t.Fatalf("split = %q, %q", serverName, toolName)
	}

	if _, _, err := SplitNamespacedName("bad"); err == nil {
		t.Fatal("SplitNamespacedName() expected error")
	}
}

type fakeSession struct {
	tools      []*mcp.Tool
	result     *mcp.CallToolResult
	err        error
	calledName string
	closed     bool
}

func (s *fakeSession) ListTools(ctx context.Context, params *mcp.ListToolsParams) (*mcp.ListToolsResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &mcp.ListToolsResult{Tools: s.tools}, nil
}

func (s *fakeSession) CallTool(ctx context.Context, params *mcp.CallToolParams) (*mcp.CallToolResult, error) {
	s.calledName = params.Name
	if s.err != nil {
		return nil, s.err
	}
	if s.result == nil {
		return nil, errors.New("missing fake result")
	}
	return s.result, nil
}

func (s *fakeSession) Close() error {
	s.closed = true
	return nil
}
