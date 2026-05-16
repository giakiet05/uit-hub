package tui

import (
	"context"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

type Runner struct {
	session       *runtime.Session
	agent         agent.Agent
	stdin         io.Reader
	stdout        io.Writer
	logs          *LogBuffer
	initialPrompt string
}

func NewRunner(session *runtime.Session, runtimeAgent agent.Agent, stdin io.Reader, stdout io.Writer, logs *LogBuffer, initialPrompt string) *Runner {
	return &Runner{
		session:       session,
		agent:         runtimeAgent,
		stdin:         stdin,
		stdout:        stdout,
		logs:          logs,
		initialPrompt: initialPrompt,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	model := newModel(ctx, r.session, r.agent, r.logs, r.initialPrompt)
	program := tea.NewProgram(
		model,
		tea.WithContext(ctx),
		tea.WithInput(r.stdin),
		tea.WithOutput(r.stdout),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := program.Run()
	return err
}
