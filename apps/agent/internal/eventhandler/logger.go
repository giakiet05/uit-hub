package eventhandler

import (
	"context"
	"log/slog"
	"sync"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eventbus"
)

// LoggerEventHandler listens to lifecycle events and writes structured logs.
type LoggerEventHandler struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	eventBus  *eventbus.EventBus
	logger    *slog.Logger
	started   bool
	eventChan eventbus.EventChan
}

// NewLoggerEventHandler creates a new logger subscriber.
func NewLoggerEventHandler(ctx context.Context, bus *eventbus.EventBus, logger *slog.Logger) *LoggerEventHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggerEventHandler{
		ctx:      ctx,
		eventBus: bus,
		logger:   logger.With("service", "logger_subscriber"),
	}
}

// Start begins listening to lifecycle events on the bus.
func (l *LoggerEventHandler) Start() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.started {
		return nil
	}

	if l.eventBus == nil {
		return nil
	}

	parentCtx := l.ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	ch := l.eventBus.Subscribe(eventbus.TopicLifecycle)
	if ch == nil {
		cancel()
		return nil
	}

	l.cancel = cancel
	l.eventChan = ch
	l.started = true

	go l.loop(ctx, ch)
	return nil
}

// Stop unsubscribes from the bus and stops the listener loop.
func (l *LoggerEventHandler) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.started {
		return
	}

	cancel := l.cancel
	ch := l.eventChan

	l.started = false
	l.cancel = nil
	l.eventChan = nil

	if cancel != nil {
		cancel()
	}

	if ch != nil && l.eventBus != nil {
		l.eventBus.Unsubscribe(eventbus.TopicLifecycle, ch)
	}
}

func (l *LoggerEventHandler) loop(ctx context.Context, ch eventbus.EventChan) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			l.handleEvent(ctx, event)
		}
	}
}

func (l *LoggerEventHandler) handleEvent(ctx context.Context, event eventbus.Event) {
	if event == nil {
		return
	}

	switch typed := event.(type) {
	case agent.RunStartedEvent:
		agent.Trace(l.logger, ctx, "agent.run.started", "Run started", "session_id", typed.SessionID, "max_rounds", typed.MaxRounds)
	case agent.RoundStartedEvent:
		agent.Trace(l.logger, ctx, "agent.round.started", "Round started", "session_id", typed.SessionID, "round", typed.Round)
	case agent.ModelCallStartedEvent:
		agent.Trace(l.logger, ctx, "agent.model.started", "Model call started", "session_id", typed.SessionID, "round", typed.Round)
	case agent.ModelCallCompletedEvent:
		agent.Trace(l.logger, ctx, "agent.model.completed", "Model call completed", "session_id", typed.SessionID, "round", typed.Round, "duration", typed.Duration)
	case agent.ToolCallStartedEvent:
		agent.Trace(l.logger, ctx, "agent.tool.started", "Tool call started", "session_id", typed.SessionID, "round", typed.Round, "call_id", typed.Call.ID, "tool_name", typed.Call.Name)
	case agent.ToolCallCompletedEvent:
		agent.Trace(l.logger, ctx, "agent.tool.completed", "Tool call completed", "session_id", typed.SessionID, "round", typed.Round, "duration", typed.Duration)
	case agent.ToolCallFailedEvent:
		agent.Trace(l.logger, ctx, "agent.tool.failed", "Tool call failed", "session_id", typed.SessionID, "round", typed.Round, "tool_name", typed.Call.Name, "duration", typed.Duration, "error", typed.Observation)
	case agent.PlanCreatedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.plan.created", "Plan created", "session_id", typed.SessionID, "steps", typed.StepCount)
	case agent.StepStartedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.step.started", "Step started", "session_id", typed.SessionID, "step_id", typed.StepID)
	case agent.StepCompletedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.step.completed", "Step completed", "session_id", typed.SessionID, "step_id", typed.StepID)
	case agent.StepFailedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.step.failed", "Step failed", "session_id", typed.SessionID, "step_id", typed.StepID, "error", typed.Err)
	case agent.ReplanStartedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.replan.started", "Replan started", "session_id", typed.SessionID, "step_id", typed.StepID)
	case agent.PlanUpdatedEvent:
		agent.Trace(l.logger, ctx, "plan_execute.plan.updated", "Plan updated", "session_id", typed.SessionID)
	case agent.FinalAnswerEvent:
		agent.Trace(l.logger, ctx, "plan_execute.finalizing", "Finalizing plan", "session_id", typed.SessionID)
	case agent.RunCompletedEvent:
		agent.LogRunStats(l.logger, "agent", ctx, typed.SessionID, &typed.Stats)
		agent.Trace(l.logger, ctx, "agent.run.completed", "Run completed", "session_id", typed.SessionID, "reason", typed.Reason)
	case agent.RunFailedEvent:
		agent.LogRunStats(l.logger, "agent", ctx, typed.SessionID, &typed.Stats)
		agent.Trace(l.logger, ctx, "agent.run.failed", "Run failed", "session_id", typed.SessionID, "error", typed.Err)
	}
}
