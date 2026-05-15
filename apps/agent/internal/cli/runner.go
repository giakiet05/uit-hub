package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

type Runner struct {
	session *runtime.Session
	loop    *agent.Loop
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	logger  *slog.Logger
}

func NewRunner(session *runtime.Session, loop *agent.Loop, stdin io.Reader, stdout io.Writer, stderr io.Writer, logger *slog.Logger) *Runner {
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	return &Runner{
		session: session,
		loop:    loop,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		logger:  logger,
	}
}

func (r *Runner) Run(ctx context.Context, args []string) error {
	scanner := bufio.NewScanner(r.stdin)
	initialPrompt := strings.TrimSpace(strings.Join(args, " "))

	r.logger.InfoContext(ctx, "Agent CLI session started", "session_id", r.session.ID, "has_initial_prompt", initialPrompt != "")
	defer r.logger.InfoContext(ctx, "Agent CLI session ended", "session_id", r.session.ID)

	if initialPrompt != "" {
		if err := r.runPrompt(ctx, initialPrompt); err != nil {
			return err
		}
	}

	for {
		if _, err := fmt.Fprint(r.stdout, "user> "); err != nil {
			return err
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				r.logger.DebugContext(ctx, "Failed to read CLI input", "session_id", r.session.ID, "error", err)
				return err
			}
			r.logger.DebugContext(ctx, "CLI input closed", "session_id", r.session.ID)
			return nil
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "" {
			r.logger.DebugContext(ctx, "Empty CLI input ignored", "session_id", r.session.ID)
			continue
		}
		if prompt == "exit" || prompt == "quit" {
			r.logger.DebugContext(ctx, "CLI exit command received", "session_id", r.session.ID)
			return nil
		}

		if err := r.runPrompt(ctx, prompt); err != nil {
			return err
		}
	}
}

func (r *Runner) runPrompt(ctx context.Context, prompt string) error {
	r.logger.DebugContext(ctx, "Running prompt through agent loop", "session_id", r.session.ID, "prompt_chars", len(prompt))
	response, err := r.loop.Run(ctx, r.session, prompt)
	if err != nil {
		r.logger.DebugContext(ctx, "Agent loop returned error", "session_id", r.session.ID, "error", err)
		return err
	}

	r.logger.DebugContext(ctx, "Writing assistant response", "session_id", r.session.ID, "response_chars", len(response.ContentText()))
	_, err = fmt.Fprintf(r.stdout, "assistant: %s\n", response.ContentText())
	return err
}
