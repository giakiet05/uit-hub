package localtool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// ReadFile reads text files from a configured sandbox root.
type ReadFile struct {
	tool.BaseTool
	rootDir string
}

// NewReadFile creates a sandboxed file-reading tool.
func NewReadFile(rootDir string) *ReadFile {
	return &ReadFile{
		BaseTool: tool.NewBaseTool(
			tool.Definition{
				Name:        "read_file",
				Description: "Read text content from a file inside the agent file sandbox.",
				InputSchema: tool.ObjectSchema(
					map[string]any{
						"path": tool.StringProperty("Relative file path inside the sandbox, for example notes/result.txt."),
					},
					"path",
				),
			},
			tool.NewReadOnlyMetadata(false),
		),
		rootDir: rootDir,
	}
}

// Execute reads a sandbox-relative file and returns the content with metadata
// as JSON.
func (t *ReadFile) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	relativePath, err := stringArg(call.Arguments, "path")
	if err != nil {
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: err}
	}

	fullPath, err := safeSandboxPath(t.rootDir, relativePath)
	if err != nil {
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("invalid path %q: %v", relativePath, err)}
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("file %q does not exist", relativePath)}
		}
		if os.IsPermission(err) {
			return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("permission denied reading %q", relativePath)}
		}
		return tool.Result{}, tool.ToolError{Type: tool.ErrorTypeValidation, Name: call.Name, Err: fmt.Errorf("could not read file %q", relativePath)}
	}

	output, err := json.Marshal(map[string]any{
		"path":    fullPath,
		"bytes":   len(data),
		"content": string(data),
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
