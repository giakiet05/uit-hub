package bus

import (
	"log/slog"
	"sync"
)

// EventChan is the channel type for receiving events.
type EventChan chan Event

// EventBus is a topic-based pub-sub event bus.
type EventBus struct {
	subscribers      map[Topic][]EventChan
	mu               sync.RWMutex
	closed           bool
	subscriberBuffer int
	logger           *slog.Logger
}

// New creates a new EventBus with the specified channel buffer size.
func New(subscriberBuffer int, logger *slog.Logger) *EventBus {
	if subscriberBuffer <= 0 {
		subscriberBuffer = 128
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &EventBus{
		subscribers:      make(map[Topic][]EventChan),
		subscriberBuffer: subscriberBuffer,
		logger:           logger.With("component", "eventbus"),
	}
}

// Close shuts down the bus and closes all subscriber channels.
func (e *EventBus) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return
	}
	e.closed = true

	for _, subscribers := range e.subscribers {
		for _, s := range subscribers {
			close(s)
		}
	}

	e.subscribers = make(map[Topic][]EventChan)
	e.logger.Debug("Event bus closed")
}

// Publish distributes an event to all subscribers of its topic.
// It uses non-blocking sends; if a subscriber's channel is full, the event is dropped.
func (e *EventBus) Publish(event Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.closed || event == nil {
		return
	}

	topic := event.Topic()
	if topic == "" {
		return
	}

	publishTo := func(targetTopic Topic) {
		if subscribers, ok := e.subscribers[targetTopic]; ok {
			for _, s := range subscribers {
				select {
				case s <- event:
				default:
					e.logger.Warn("Event dropped - subscriber channel full", "topic", string(targetTopic))
				}
			}
		}
	}

	publishTo(topic)
	publishTo(TopicAll)
}

// Subscribe registers a new channel to receive events for a specific topic.
func (e *EventBus) Subscribe(topic Topic) EventChan {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed || topic == "" {
		return nil
	}

	ch := make(EventChan, e.subscriberBuffer)

	e.subscribers[topic] = append(e.subscribers[topic], ch)
	totalSubscribers := len(e.subscribers[topic])
	e.logger.Debug("Subscriber registered", "topic", string(topic), "total_subscribers", totalSubscribers)

	return ch
}

// Unsubscribe removes a channel from a topic and closes it.
func (e *EventBus) Unsubscribe(topic Topic, ch EventChan) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed || topic == "" || ch == nil {
		return
	}

	subscribers, ok := e.subscribers[topic]
	if !ok || len(subscribers) == 0 {
		return
	}

	filtered := subscribers[:0]
	for _, sub := range subscribers {
		if sub == ch {
			close(ch)
		} else {
			filtered = append(filtered, sub)
		}
	}

	for i := len(filtered); i < len(subscribers); i++ {
		subscribers[i] = nil
	}

	if len(filtered) == 0 {
		delete(e.subscribers, topic)
		e.logger.Debug("Subscriber unregistered", "topic", string(topic), "total_subscribers", 0)
		return
	}

	e.subscribers[topic] = filtered
	e.logger.Debug("Subscriber unregistered", "topic", string(topic), "total_subscribers", len(filtered))
}
