package localtool

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestWriteFileWritesInsideSandbox(t *testing.T) {
	rootDir := t.TempDir()
	fileTool := NewWriteFile(rootDir)

	result, err := fileTool.Execute(context.Background(), tool.Call{
		ID:   "call-write",
		Name: "write_file",
		Arguments: map[string]any{
			"path":    "notes/result.txt",
			"content": "hello file",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, `"bytes":10`) {
		t.Fatalf("result content = %q, want byte count", result.Content)
	}

	data, err := os.ReadFile(filepath.Join(rootDir, "notes", "result.txt"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got, want := string(data), "hello file"; got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestWriteFileRejectsPathTraversal(t *testing.T) {
	fileTool := NewWriteFile(t.TempDir())

	_, err := fileTool.Execute(context.Background(), tool.Call{
		ID:   "call-write",
		Name: "write_file",
		Arguments: map[string]any{
			"path":    "../escape.txt",
			"content": "bad",
		},
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want path traversal error")
	}
}
