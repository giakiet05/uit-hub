package main

import (
	"fmt"
	"os"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/academic/internal/client"
	"github.com/giakiet05/uit-hub/apps/mcp-servers/academic/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"uit-hub-mcp-academic",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	c := client.New()

	tools.RegisterAuth(s, c)
	tools.RegisterAcademic(s, c)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "mcp server error: %v\n", err)
		os.Exit(1)
	}
}
