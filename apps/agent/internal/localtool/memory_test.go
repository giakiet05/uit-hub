package localtool

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestMemoryToolsWriteListAndRead(t *testing.T) {
	store := memory.NewFileStore(t.TempDir())
	writeTool := NewMemoryWrite(store)
	listTool := NewMemoryList(store)
	readTool := NewMemoryRead(store)

	writeResult, err := writeTool.Execute(context.Background(), tool.Call{
		ID:   "call-write",
		Name: "memory_write",
		Arguments: map[string]any{
			"type":        "feedback",
			"name":        "Response Style",
			"description": "User prefers concise answers.",
			"content":     "Keep answers short and direct.",
		},
	})
	if err != nil {
		t.Fatalf("memory_write Execute() error = %v", err)
	}
	if !strings.Contains(writeResult.Content, `"id":"feedback_response_style"`) {
		t.Fatalf("write result = %q", writeResult.Content)
	}

	listResult, err := listTool.Execute(context.Background(), tool.Call{
		ID:        "call-list",
		Name:      "memory_list",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("memory_list Execute() error = %v", err)
	}
	if !strings.Contains(listResult.Content, "feedback_response_style") {
		t.Fatalf("list result = %q", listResult.Content)
	}

	readResult, err := readTool.Execute(context.Background(), tool.Call{
		ID:   "call-read",
		Name: "memory_read",
		Arguments: map[string]any{
			"id": "feedback_response_style",
		},
	})
	if err != nil {
		t.Fatalf("memory_read Execute() error = %v", err)
	}
	if !strings.Contains(readResult.Content, "Keep answers short and direct.") {
		t.Fatalf("read result = %q", readResult.Content)
	}
}
