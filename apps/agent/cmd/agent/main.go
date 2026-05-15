package main

import (
	"context"
	"fmt"
	"os"

	"github.com/giakiet05/uit-hub/apps/agent/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
