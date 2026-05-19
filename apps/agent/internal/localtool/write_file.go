package localtool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// WriteFile writes files under a configured sandbox root.
type WriteFile struct {
	rootDir string
}

// NewWriteFile creates a sandboxed file-writing tool.
func NewWriteFile(rootDir string) *WriteFile {
	return &WriteFile{rootDir: rootDir}
}

// Definition describes the sandboxed file-writing tool schema.
func (t *WriteFile) Definition() tool.Definition {
	return tool.Definition{
		Name:        "write_file",
		Description: "Write text content to a file inside the agent file sandbox.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Relative file path inside the sandbox, for example notes/result.txt.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "Text content to write.",
				},
			},
			"required": []string{"path", "content"},
		},
	}
}

// Execute writes text content to a sandbox-relative path and returns metadata as
// JSON.
func (t *WriteFile) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	relativePath, err := stringArg(call.Arguments, "path")
	if err != nil {
		return tool.Result{}, err
	}
	content, err := stringArg(call.Arguments, "content")
	if err != nil {
		return tool.Result{}, err
	}

	fullPath, err := t.safePath(relativePath)
	if err != nil {
		return tool.Result{}, err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return tool.Result{}, fmt.Errorf("create parent directory: %w", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return tool.Result{}, fmt.Errorf("write file: %w", err)
	}

	output, err := json.Marshal(map[string]any{
		"path":  fullPath,
		"bytes": len([]byte(content)),
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal result: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}

// safePath resolves a requested relative path and rejects attempts to escape the
// sandbox root.
func (t *WriteFile) safePath(relativePath string) (string, error) {
	if strings.TrimSpace(relativePath) == "" {
		return "", errors.New("path must not be empty")
	}
	if filepath.IsAbs(relativePath) {
		return "", errors.New("path must be relative")
	}

	cleanPath := filepath.Clean(relativePath)
	if cleanPath == "." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) || cleanPath == ".." {
		return "", errors.New("path escapes sandbox")
	}

	rootDir := t.rootDir
	if rootDir == "" {
		rootDir = filepath.Join("tmp", "agent-files")
	}
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("resolve sandbox root: %w", err)
	}

	fullPath := filepath.Join(rootDir, cleanPath)
	if !strings.HasPrefix(fullPath, rootDir+string(filepath.Separator)) && fullPath != rootDir {
		return "", errors.New("path escapes sandbox")
	}
	return fullPath, nil
}
