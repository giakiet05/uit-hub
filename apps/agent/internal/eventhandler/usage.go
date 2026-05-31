package eventhandler

import (
	"context"
	"log/slog"
	"sync"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eventbus"
)

// UsageTracker is an interface for tracking usage statistics.
type UsageTracker interface {
	AddRunStats(stats agent.RunStats)
}

// UsageEventHandler listens to agent events to record token usage and run stats.
type UsageEventHandler struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	eventBus  *eventbus.EventBus
	logger    *slog.Logger
	started   bool
	eventChan eventbus.EventChan
	tracker   UsageTracker
}

// NewUsageEventHandler creates a new usage metric event handler.
func NewUsageEventHandler(ctx context.Context, tracker UsageTracker, bus *eventbus.EventBus, logger *slog.Logger) *UsageEventHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &UsageEventHandler{
		ctx:      ctx,
		tracker:  tracker,
		eventBus: bus,
		logger:   logger.With("service", "usage_event_handler"),
	}
}

// Start begins listening to lifecycle events on the bus.
func (u *UsageEventHandler) Start() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.started {
		return nil
	}

	if u.eventBus == nil || u.tracker == nil {
		return nil
	}

	parentCtx := u.ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	ch := u.eventBus.Subscribe(eventbus.TopicLifecycle)
	if ch == nil {
		cancel()
		return nil
	}

	u.cancel = cancel
	u.eventChan = ch
	u.started = true

	u.logger.Debug("Usage subscriber started")
	go u.loop(ctx, ch)
	return nil
}

// Stop unsubscribes from the bus and stops the listener loop.
func (u *UsageEventHandler) Stop() {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.started {
		return
	}

	cancel := u.cancel
	ch := u.eventChan

	u.started = false
	u.cancel = nil
	u.eventChan = nil

	if cancel != nil {
		cancel()
	}

	if ch != nil && u.eventBus != nil {
		u.eventBus.Unsubscribe(eventbus.TopicLifecycle, ch)
	}

	u.logger.Debug("Usage subscriber stopped")
}

func (u *UsageEventHandler) loop(ctx context.Context, ch eventbus.EventChan) {
	for {
		select {
		case <-ctx.Done():
			u.logger.Debug("Usage subscriber loop exiting")
			return
		case event, ok := <-ch:
			if !ok {
				u.logger.Debug("Usage subscriber event channel closed")
				return
			}
			u.handleEvent(event)
		}
	}
}

func (u *UsageEventHandler) handleEvent(event eventbus.Event) {
	if event == nil {
		return
	}

	switch typed := event.(type) {
	case agent.RunCompletedEvent:
		u.addRunStats(typed.Stats)
	case agent.RunFailedEvent:
		u.addRunStats(typed.Stats)
	}
}

func (u *UsageEventHandler) addRunStats(stats agent.RunStats) {
	u.tracker.AddRunStats(stats)
}
