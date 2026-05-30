// Package tui provides the terminal UI for interactive agent sessions.
package tui

import (
	"context"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

// Runner owns the Bubble Tea program setup for the agent TUI.
type Runner struct {
	session       *session.Session
	stdin         io.Reader
	stdout        io.Writer
	logs          *LogBuffer
	initialPrompt string
}

// NewRunner creates a TUI runner for one session.
func NewRunner(session *session.Session, stdin io.Reader, stdout io.Writer, logs *LogBuffer, initialPrompt string) *Runner {
	return &Runner{
		session:       session,
		stdin:         stdin,
		stdout:        stdout,
		logs:          logs,
		initialPrompt: initialPrompt,
	}
}

// Run starts the Bubble Tea event loop.
func (r *Runner) Run(ctx context.Context) error {
	model := newModel(ctx, r.session, r.logs, r.initialPrompt)
	program := tea.NewProgram(
		model,
		tea.WithContext(ctx),
		tea.WithInput(r.stdin),
		tea.WithOutput(r.stdout),
		tea.WithAltScreen(),
	)
	_, err := program.Run()
	return err
}
