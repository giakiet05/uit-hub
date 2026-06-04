package app

import (
	"context"
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/event/bus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/event/handler"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

// CompositeCloser bundles multiple io.Closer implementations into one.
type CompositeCloser struct {
	closers []io.Closer
}

// Close calls Close on all bundled closers.
func (c *CompositeCloser) Close() error {
	var firstErr error
	for _, closer := range c.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// StartEventHandlers initializes and starts all necessary event handlers.
// It returns an io.Closer that stops all handlers when called.
func StartEventHandlers(ctx context.Context, bus *bus.EventBus, sessionState *session.State, db *storage.DB, logger *slog.Logger) io.Closer {
	var closers []io.Closer

	usageSub := handler.NewUsageEventHandler(ctx, sessionState, bus, logger)
	usageSub.Start()
	// UsageEventHandler needs a Stop method, which we can wrap in a closer
	closers = append(closers, &handlerCloser{stop: usageSub.Stop})

	loggerSub := handler.NewLoggerEventHandler(ctx, bus, logger)
	loggerSub.Start()
	closers = append(closers, &handlerCloser{stop: loggerSub.Stop})

	if db != nil {
		storageSub := handler.NewStorageEventHandler(ctx, db, sessionState, bus, logger)
		storageSub.Start()
		closers = append(closers, &handlerCloser{stop: storageSub.Stop})
	}

	return &CompositeCloser{closers: closers}
}

type handlerCloser struct {
	stop func()
}

func (hc *handlerCloser) Close() error {
	hc.stop()
	return nil
}
