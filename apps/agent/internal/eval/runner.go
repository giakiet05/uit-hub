package eval

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/giakiet05/uit-hub/apps/agent/internal/app"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
)

// Runner executes eval cases against the shared agent runtime.
type Runner struct {
	cfg    config.Config
	logger *slog.Logger
}

// NewRunner creates an eval runner.
func NewRunner(cfg config.Config, logger *slog.Logger) *Runner {
	return &Runner{
		cfg:    cfg,
		logger: logger,
	}
}

// RunCase runs one eval case and writes a report.
func (r *Runner) RunCase(ctx context.Context, testCase Case, outputRoot string) (Result, string, error) {
	cfg := r.cfg
	cfg.Memory.Path = filepath.Join(r.cfg.Memory.Path, safePathPart(testCase.ID))

	runtime, err := app.NewRuntime(ctx, cfg, r.logger)
	if err != nil {
		return Result{}, "", err
	}
	defer runtime.Close(ctx)

	recorder := NewRecorder(testCase)
	events := runtime.Agent.Run(ctx, runtime.Session, testCase.Prompt)
	recorder.Consume(ctx, events)

	result := recorder.Finish(Judge(testCase, recorder.result))
	reportDir, err := WriteReport(outputRoot, testCase, result, recorder.Trace())
	if err != nil {
		return Result{}, "", err
	}
	return result, reportDir, nil
}
