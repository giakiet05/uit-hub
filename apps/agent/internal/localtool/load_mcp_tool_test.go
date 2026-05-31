package localtool

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestLoadMCPToolRegistersRuntimeTool(t *testing.T) {
	runtime := tool.NewRuntimeRegistry(2)
	loader := fakeMCPToolLoader{
		tools: map[string]tool.Tool{
			"uit__get_student": fakeTool{name: "uit__get_student"},
		},
	}
	loadTool := NewLoadMCPTool(loader, runtime)

	result, err := loadTool.Execute(context.Background(), tool.Call{
		ID:   "call-1",
		Name: "load_mcp_tool",
		Arguments: map[string]any{
			"name": "uit__get_student",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !runtime.Has("uit__get_student") {
		t.Fatal("runtime registry missing loaded tool")
	}
	if !strings.Contains(result.Content, "loaded") {
		t.Fatalf("result content = %q, want loaded message", result.Content)
	}
}

func TestLoadMCPToolAlreadyLoaded(t *testing.T) {
	runtime := tool.NewRuntimeRegistry(2)
	if err := runtime.Register(fakeTool{name: "uit__get_student"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	loadTool := NewLoadMCPTool(fakeMCPToolLoader{}, runtime)

	result, err := loadTool.Execute(context.Background(), tool.Call{
		ID:   "call-1",
		Name: "load_mcp_tool",
		Arguments: map[string]any{
			"name": "uit__get_student",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, "already loaded") {
		t.Fatalf("result content = %q, want already loaded message", result.Content)
	}
}

func TestLoadMCPToolUnknown(t *testing.T) {
	loadTool := NewLoadMCPTool(fakeMCPToolLoader{}, tool.NewRuntimeRegistry(2))

	_, err := loadTool.Execute(context.Background(), tool.Call{
		ID:   "call-1",
		Name: "load_mcp_tool",
		Arguments: map[string]any{
			"name": "missing__tool",
		},
	})
	if err == nil {
		t.Fatal("Execute() expected error")
	}
	if !strings.Contains(err.Error(), "unknown MCP tool") {
		t.Fatalf("error message = %q, want unknown tool message", err.Error())
	}
}

type fakeMCPToolLoader struct {
	tools map[string]tool.Tool
}

func (l fakeMCPToolLoader) LoadTool(name string) (tool.Tool, bool) {
	loadedTool, exists := l.tools[name]
	return loadedTool, exists
}

type fakeTool struct {
	name string
}

func (t fakeTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        t.name,
		Description: "fake tool",
		InputSchema: tool.EmptyInputSchema(),
	}
}

func (t fakeTool) Metadata() tool.Metadata {
	return tool.NewReadOnlyMetadata(true, tool.DefaultMaxResultChars)
}

func (t fakeTool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}
	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: "ok",
	}, nil
}
