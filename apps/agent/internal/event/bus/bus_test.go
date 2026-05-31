package bus_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
)

type mockEvent struct {
	topic bus.Topic
	data  string
}

func (m mockEvent) Topic() bus.Topic {
	return m.topic
}

func TestEventBus(t *testing.T) {
	logger := slog.Default()
	bus := bus.New(10, logger)

	ch := bus.Subscribe("test_topic")
	if ch == nil {
		t.Fatal("Expected channel, got nil")
	}

	bus.Publish(mockEvent{topic: "test_topic", data: "hello"})

	select {
	case e := <-ch:
		if e.(mockEvent).data != "hello" {
			t.Errorf("Expected 'hello', got '%s'", e.(mockEvent).data)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for event")
	}

	// Test Unsubscribe
	bus.Unsubscribe("test_topic", ch)
	
	// Should not block or panic when publishing to unsubscribed topic
	bus.Publish(mockEvent{topic: "test_topic", data: "hello2"})
	
	// Test TopicAll
	chAll := bus.Subscribe(bus.TopicAll)
	bus.Publish(mockEvent{topic: "other_topic", data: "all"})
	
	select {
	case e := <-chAll:
		if e.(mockEvent).data != "all" {
			t.Errorf("Expected 'all', got '%s'", e.(mockEvent).data)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for wildcard event")
	}

	bus.Close()
	
	// Should not panic after close
	bus.Publish(mockEvent{topic: "test_topic", data: "after_close"})
}
