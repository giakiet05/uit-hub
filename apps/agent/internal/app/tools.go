package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// newToolSet registers the base tools used by the current baseline agent and
// creates an empty runtime registry for tools discovered during one run.
func newToolSet(
	ctx context.Context,
	cfg config.Config,
	logger *slog.Logger,
	memoryStore memory.Store,
) (*tool.ToolSet, []io.Closer, error) {
	baseRegistry, err := tool.NewBaseRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
		localtool.NewCurrentTime(),
		localtool.NewReadFile("tmp/agent-files"),
		localtool.NewWriteFile("tmp/agent-files"),
		localtool.NewMemoryRead(memoryStore),
		localtool.NewMemoryWrite(memoryStore),
		localtool.NewMemoryList(memoryStore),
	)
	if err != nil {
		return nil, nil, err
	}

	runtimeRegistry := tool.NewRuntimeRegistry(cfg.Tool.RuntimeLimit)

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
			if err := baseRegistry.Register(mcpTool); err != nil {
				return nil, closers, err
			}
		}
		logger.DebugContext(ctx, "MCP tools registered", "server_name", server.Name, "tool_count", len(mcpTools))
	}

	return tool.NewToolSet(baseRegistry, runtimeRegistry), closers, nil
}

// closeAll attempts to close all MCP clients and logs any errors without interrupting the closing of other
func closeAll(ctx context.Context, logger *slog.Logger, closers []io.Closer) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			logger.DebugContext(ctx, "Close resource failed", "error", err)
		}
	}
}

type mcpServerLoadReport struct {
	ServerName string
	Client     *mcpadapter.Client
	ToolCount  int
	Success    bool
	Attempts   int
	Err        error
}

func loadMCPServers(ctx context.Context, logger *slog.Logger, cfg config.Config) ([]mcpServerLoadReport, error) {
	reportCh := make(chan mcpServerLoadReport, len(cfg.MCP.Servers))
	var wg sync.WaitGroup

	for _, server := range cfg.MCP.Servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reportCh <- loadOneMCPServerWithRetry(ctx, server, cfg)
		}()
	}

	go func() {
		wg.Wait()
		close(reportCh)
	}()

	reports := make([]mcpServerLoadReport, len(cfg.MCP.Servers))
	for r := range reportCh {
		reports = append(reports, r)
	}

	return reports, nil
}

func loadOneMCPServerWithRetry(ctx context.Context, server config.MCPServerConfig, cfg config.Config, registry) mcpServerLoadReport {
	attemptCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	for i := 1; i <= 3; i++ {
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
			continue
		}

		mcpTools, err := client.Tools(ctx)
		if err != nil {
			continue
		}
		for _, mcpTool := range mcpTools {
			if err := baseRegistry.Register(mcpTool); err != nil {
				continue
			}
		}

	}

}
