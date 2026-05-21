package localtool

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestReadFileReadsInsideSandbox(t *testing.T) {
	rootDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(rootDir, "notes"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(rootDir, "notes", "result.txt"), []byte("hello file"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	fileTool := NewReadFile(rootDir)
	result, err := fileTool.Execute(context.Background(), tool.Call{
		ID:   "call-read",
		Name: "read_file",
		Arguments: map[string]any{
			"path": "notes/result.txt",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, `"bytes":10`) {
		t.Fatalf("result content = %q, want byte count", result.Content)
	}
	if !strings.Contains(result.Content, `"content":"hello file"`) {
		t.Fatalf("result content = %q, want file content", result.Content)
	}
}

func TestReadFileRejectsPathTraversal(t *testing.T) {
	fileTool := NewReadFile(t.TempDir())

	_, err := fileTool.Execute(context.Background(), tool.Call{
		ID:   "call-read",
		Name: "read_file",
		Arguments: map[string]any{
			"path": "../escape.txt",
		},
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want path traversal error")
	}
}
