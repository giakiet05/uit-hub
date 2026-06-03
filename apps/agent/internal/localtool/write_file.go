package localtool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// WriteFile writes files under a configured sandbox root.
type WriteFile struct {
	tool.BaseTool
	rootDir string
}

// NewWriteFile creates a sandboxed file-writing tool.
func NewWriteFile(rootDir string) *WriteFile {
	return &WriteFile{
		BaseTool: tool.NewBaseTool(
			tool.Definition{
				Name:        "write_file",
				Description: "Write text content to a file inside the agent file sandbox.",
				InputSchema: tool.ObjectSchema(
					map[string]any{
						"path":    tool.StringProperty("Relative file path inside the sandbox, for example notes/result.txt."),
						"content": tool.StringProperty("Text content to write."),
					},
					"path",
					"content",
				),
			},
			tool.NewWriteMetadata(false, true),
		),
		rootDir: rootDir,
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
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: err}
	}
	content, err := stringArg(call.Arguments, "content")
	if err != nil {
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: err}
	}

	fullPath, err := safeSandboxPath(t.rootDir, relativePath)
	if err != nil {
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("invalid path %q: %v", relativePath, err)}
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		if os.IsPermission(err) {
			return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("permission denied creating directory for %q", relativePath)}
		}
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("could not create parent directory for %q", relativePath)}
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		if os.IsPermission(err) {
			return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("permission denied writing to %q", relativePath)}
		}
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("could not write file %q", relativePath)}
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
