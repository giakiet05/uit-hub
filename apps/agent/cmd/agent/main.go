package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/giakiet05/uit-hub/apps/agent/internal/app"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
)

// main starts the agent TUI and reports startup/runtime errors to stderr.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Mock SSO Login step for Prototype:
	// In a real app, this token would come from a login prompt or OAuth flow.
	ctx = context.WithValue(ctx, mcpadapter.TokenKey, "mock-22520001")

	if err := app.Run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
