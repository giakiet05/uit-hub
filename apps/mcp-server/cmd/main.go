package main

import (
	"fmt"
	"os"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"uit-hub-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	c := client.New()

	// Register all tool groups
	tools.RegisterAuth(s, c)
	tools.RegisterStudent(s, c)
	tools.RegisterCourse(s, c)
	tools.RegisterCTSV(s, c)
	tools.RegisterRoom(s, c)

	// Serve over stdio (standard MCP transport)
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "mcp server error: %v\n", err)
		os.Exit(1)
	}
}
