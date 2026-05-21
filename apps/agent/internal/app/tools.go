package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// newToolRegistry registers the temporary local tools used by the current
// baseline agent.
func newToolRegistry(ctx context.Context, cfg config.Config, logger *slog.Logger) (*tool.Registry, []io.Closer, error) {
	registry, err := tool.NewRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
		localtool.NewCurrentTime(),
		localtool.NewReadFile("tmp/agent-files"),
		localtool.NewWriteFile("tmp/agent-files"),
	)
	if err != nil {
		return nil, nil, err
	}

	var closers []io.Closer
	for _, server := range cfg.MCP.Servers {
		if server.Transport != "stdio" {
			return nil, closers, fmt.Errorf("unsupported MCP transport %q for server %q", server.Transport, server.Name)
		}
		client, err := mcpadapter.NewStdioClient(ctx, mcpadapter.StdioConfig{
			ServerName: server.Name,
			Command:    server.Command,
			Args:       server.Args,
			WorkDir:    server.WorkDir,
			Client: mcpadapter.ClientIdentity{
				Name:    cfg.MCP.ClientName,
				Version: cfg.MCP.ClientVersion,
			},
		})
		if err != nil {
			return nil, nil, err
		}
		closers = append(closers, client)

		mcpTools, err := client.Tools(ctx)
		if err != nil {
			return nil, closers, err
		}
		for _, mcpTool := range mcpTools {
			if err := registry.Register(mcpTool); err != nil {
				return nil, closers, err
			}
		}
		logger.DebugContext(ctx, "MCP tools registered", "server_name", server.Name, "tool_count", len(mcpTools))
	}

	return registry, closers, nil
}

func closeAll(ctx context.Context, logger *slog.Logger, closers []io.Closer) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			logger.DebugContext(ctx, "Close resource failed", "error", err)
		}
	}
}
