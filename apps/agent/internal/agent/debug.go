package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

// WritePromptDebugFile writes the latest LLM message stack to tmp/prompt.md for
// local inspection without exposing it through the TUI log stream.
func WritePromptDebugFile(ctx context.Context, logger *slog.Logger, messages []conversation.Message) {
	if err := os.MkdirAll("tmp", 0o755); err != nil {
		Trace(logger, ctx, "agent.prompt_debug.write_failed", "Prompt debug directory creation failed", "error", err)
		return
	}
	if err := os.WriteFile(filepath.Join("tmp", "prompt.md"), []byte(renderPromptDebug(messages)), 0o600); err != nil {
		Trace(logger, ctx, "agent.prompt_debug.write_failed", "Prompt debug file write failed", "error", err)
	}
}

func renderPromptDebug(messages []conversation.Message) string {
	parts := make([]string, 0, len(messages))
	for _, message := range messages {
		role := conversation.RoleOf(message)
		text := conversation.Text(message)
		if text == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("## %s\n\n%s", role, text))
	}
	return strings.Join(parts, "\n\n")
}
