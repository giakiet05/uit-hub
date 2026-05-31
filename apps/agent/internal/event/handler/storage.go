package handler

import (
	"context"
	"log/slog"
	"sync"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

// StorageEventHandler listens to lifecycle events to persist session history to DB.
type StorageEventHandler struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	eventBus  *bus.EventBus
	logger    *slog.Logger
	started   bool
	eventChan bus.EventChan
	db        *storage.DB
	state     *session.State // Needed to read conversation & usage on run completion
}

// NewStorageEventHandler creates a new storage event handler.
func NewStorageEventHandler(ctx context.Context, db *storage.DB, state *session.State, eventBus *bus.EventBus, logger *slog.Logger) *StorageEventHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &StorageEventHandler{
		ctx:      ctx,
		db:       db,
		state:    state,
		eventBus: eventBus,
		logger:   logger.With("service", "storage_event_handler"),
	}
}

// Start begins listening to lifecycle events on the bus.
func (h *StorageEventHandler) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.started {
		return nil
	}

	if h.eventBus == nil || h.db == nil || h.state == nil {
		return nil
	}

	parentCtx := h.ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	ch := h.eventBus.Subscribe(bus.TopicLifecycle)
	if ch == nil {
		cancel()
		return nil
	}

	h.cancel = cancel
	h.eventChan = ch
	h.started = true

	h.logger.Debug("Storage event handler started")
	go h.loop(ctx, ch)
	return nil
}

// Stop unsubscribes from the bus and stops the listener loop.
func (h *StorageEventHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.started {
		return
	}

	cancel := h.cancel
	ch := h.eventChan

	h.started = false
	h.cancel = nil
	h.eventChan = nil

	if cancel != nil {
		cancel()
	}

	if ch != nil && h.eventBus != nil {
		h.eventBus.Unsubscribe(bus.TopicLifecycle, ch)
	}

	h.logger.Debug("Storage event handler stopped")
}

func (h *StorageEventHandler) loop(ctx context.Context, ch bus.EventChan) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("Storage event handler loop exiting")
			return
		case event, ok := <-ch:
			if !ok {
				h.logger.Debug("Storage event handler event channel closed")
				return
			}
			h.handleEvent(event)
		}
	}
}

func (h *StorageEventHandler) handleEvent(event bus.Event) {
	if event == nil {
		return
	}

	// We only care about terminal events to save state
	switch event.(type) {
	case agent.RunCompletedEvent, agent.RunFailedEvent:
		h.saveSession()
	}
}

func (h *StorageEventHandler) saveSession() {
	id := h.state.ID
	startedAt := h.state.StartedAt
	messages := h.state.Conversation.Messages()
	usage := h.state.GetUsage()

	err := h.db.UpsertSession(id, startedAt, messages, usage)
	if err != nil {
		h.logger.Error("Failed to upsert session to storage", "session_id", id, "error", err)
	} else {
		h.logger.Debug("Session saved to storage", "session_id", id, "message_count", len(messages))
	}
}
