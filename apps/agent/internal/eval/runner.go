package eval

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/planexecute"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/react"
	"github.com/giakiet05/uit-hub/apps/agent/internal/app"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
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
	if testCase.AgentType != "" {
		cfg.Agent.Type = agent.Type(testCase.AgentType)
	}

	runtime, err := app.NewRuntime(ctx, cfg, r.logger, false, "")
	if err != nil {
		return Result{}, "", err
	}
	defer runtime.Close(ctx)

	runtimeAgent, err := newEvalAgent(cfg)
	if err != nil {
		return Result{}, "", err
	}
	runtime.Session.SetAgent(runtimeAgent)

	recorder := NewRecorder(testCase)
	eventChan := runtime.Bus.Subscribe(bus.TopicAll)
	events := agentEvents(eventChan)
	go runtime.Session.Run(ctx, testCase.Prompt)
	recorder.Consume(ctx, events)

	result := recorder.Finish(Judge(testCase, recorder.result))
	reportDir, err := WriteReport(outputRoot, testCase, result, recorder.Trace())
	if err != nil {
		return Result{}, "", err
	}
	return result, reportDir, nil
}

func newEvalAgent(cfg config.Config) (agent.Agent, error) {
	providerName := os.Getenv("DEFAULT_PROVIDER")
	if providerName == "" {
		providerName = "openai"
	}
	modelName := os.Getenv("DEFAULT_MODEL")
	if modelName == "" {
		modelName = "gpt-5.4-mini"
	}
	model, err := cfg.Registry.GetModel(providerName, modelName)
	if err != nil {
		return nil, err
	}

	if cfg.Agent.Type == agent.TypePlanAndExecute {
		return planexecute.NewPlanAndExecuteAgent(
			model,
			planexecute.WithPlanMaxSteps(cfg.Agent.MaxRounds),
			planexecute.WithPlanToolTimeout(cfg.Agent.ToolTimeout),
		), nil
	}

	return react.NewReActAgent(
		model,
		react.WithMaxRounds(cfg.Agent.MaxRounds),
		react.WithToolTimeout(cfg.Agent.ToolTimeout),
		react.WithCompaction(conversation.CompactionOptions{
			DisableSnip:                       !cfg.Compaction.SnipEnabled,
			DisableMicrocompact:               !cfg.Compaction.MicrocompactEnabled,
			MicrocompactMinChars:              cfg.Compaction.MicrocompactMinChars,
			MicrocompactKeepRecentToolResults: cfg.Compaction.MicrocompactKeepRecentToolResults,
			MicrocompactTriggerRatio:          cfg.Compaction.MicrocompactTriggerRatio,
			SnipTriggerRatio:                  cfg.Compaction.SnipTriggerRatio,
		}),
		react.WithContextCollapse(react.ContextCollapseOptions{
			Enabled:          cfg.Compaction.ContextCollapseEnabled,
			TriggerRatio:     cfg.Compaction.ContextCollapseTriggerRatio,
			TargetRatio:      cfg.Compaction.ContextCollapseTargetRatio,
			KeepRecentTurns:  cfg.Compaction.ContextCollapseKeepRecentTurns,
			MaxSegments:      cfg.Compaction.ContextCollapseMaxSegments,
			MaxConcurrency:   cfg.Compaction.ContextCollapseMaxConcurrency,
			MinSegmentTokens: cfg.Compaction.ContextCollapseMinSegmentTokens,
		}),
	), nil
}

func agentEvents(eventChan bus.EventChan) <-chan agent.Event {
	events := make(chan agent.Event)
	go func() {
		defer close(events)
		for event := range eventChan {
			agentEvent, ok := event.(agent.Event)
			if !ok {
				continue
			}
			events <- agentEvent
			switch agentEvent.(type) {
			case agent.RunCompletedEvent, agent.RunFailedEvent:
				return
			}
		}
	}()
	return events
}
