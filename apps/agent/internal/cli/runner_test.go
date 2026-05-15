package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/echo"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

func TestRunnerUsesArgsAsInitialPrompt(t *testing.T) {
	var stdout bytes.Buffer
	runner := NewRunner(
		runtime.NewSession(),
		agent.NewLoop(echo.NewProvider(), prompt.NewBuilder("system"), logging.NewNopLogger()),
		strings.NewReader("exit\n"),
		&stdout,
		&bytes.Buffer{},
		logging.NewNopLogger(),
	)

	err := runner.Run(context.Background(), []string{"hello", "agent"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := stdout.String(), "assistant: received: hello agent\nuser> "; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunnerRunsInteractivePrompts(t *testing.T) {
	var stdout bytes.Buffer
	runner := NewRunner(
		runtime.NewSession(),
		agent.NewLoop(echo.NewProvider(), prompt.NewBuilder("system"), logging.NewNopLogger()),
		strings.NewReader("hello from stdin\nexit\n"),
		&stdout,
		&bytes.Buffer{},
		logging.NewNopLogger(),
	)

	err := runner.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := stdout.String(), "user> assistant: received: hello from stdin\nuser> "; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
