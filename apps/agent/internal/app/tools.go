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

const (
	mcpLoadMaxAttempts    = 3
	mcpLoadAttemptTimeout = 3 * time.Second
)

// newToolSet registers base tools and prepares deferred MCP tool loading.
func newToolSet(
	ctx context.Context,
	cfg config.Config,
	logger *slog.Logger,
	memoryStore memory.Store,
) (*tool.ToolSet, *mcpadapter.Manager, []io.Closer, error) {
	runtimeRegistry := tool.NewRuntimeRegistry(cfg.Tool.RuntimeLimit)
	mcpManager := mcpadapter.NewManager()

	baseRegistry, err := tool.NewBaseRegistry(
		localtool.NewCalculator(),
		localtool.NewReadFile("tmp/agent-files"),
		localtool.NewWriteFile("tmp/agent-files"),
		localtool.NewMemoryRead(memoryStore),
		localtool.NewMemoryWrite(memoryStore),
		localtool.NewMemoryList(memoryStore),
		localtool.NewLoadMCPTool(mcpManager, runtimeRegistry),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	reports := loadMCPServers(ctx, cfg)
	var loadedServers int
	var failedServers int
	var catalogTools int
	for _, report := range reports {
		if !report.Success {
			failedServers++
			logger.DebugContext(
				ctx,
				"MCP server load failed",
				"server_name", report.ServerName,
				"attempts", report.Attempts,
				"error", report.Err,
			)
			continue
		}

		if err := mcpManager.AddLoadedTools(report.ServerName, report.ServerDescription, report.Client, report.Tools); err != nil {
			failedServers++
			closeMCPClient(ctx, logger, report.Client)
			logger.DebugContext(
				ctx,
				"MCP server catalog registration failed",
				"server_name", report.ServerName,
				"tool_count", len(report.Tools),
				"error", err,
			)
			continue
		}

		loadedServers++
		catalogTools += len(report.Tools)
		logger.DebugContext(
			ctx,
			"MCP server loaded",
			"server_name", report.ServerName,
			"attempts", report.Attempts,
			"tool_count", len(report.Tools),
		)
	}
	logger.InfoContext(
		ctx,
		"MCP load summary",
		"configured_servers", len(cfg.MCP.Servers),
		"loaded_servers", loadedServers,
		"failed_servers", failedServers,
		"catalog_tools", catalogTools,
	)

	closers := []io.Closer{mcpManager}
	return tool.NewToolSet(baseRegistry, runtimeRegistry), mcpManager, closers, nil
}

func closeAll(ctx context.Context, logger *slog.Logger, closers []io.Closer) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			logger.DebugContext(ctx, "Close resource failed", "error", err)
		}
	}
}

type mcpServerLoadReport struct {
	ServerName        string
	ServerDescription string
	Index             int
	Client            *mcpadapter.Client
	Tools             []tool.Tool
	Success           bool
	Attempts          int
	Err               error
}

func loadMCPServers(ctx context.Context, cfg config.Config) []mcpServerLoadReport {
	if len(cfg.MCP.Servers) == 0 {
		return nil
	}

	loadID := fmt.Sprintf("%d", time.Now().UnixNano())
	reportCh := make(chan mcpServerLoadReport, len(cfg.MCP.Servers))
	var wg sync.WaitGroup
	for index, server := range cfg.MCP.Servers {
		wg.Add(1)
		go func(index int, server config.MCPServerConfig) {
			defer wg.Done()
			reportCh <- loadOneMCPServerWithRetry(ctx, index, server, cfg.MCP, loadID)
		}(index, server)
	}

	go func() {
		wg.Wait()
		close(reportCh)
	}()

	reports := make([]mcpServerLoadReport, len(cfg.MCP.Servers))
	for report := range reportCh {
		reports[report.Index] = report
	}
	return reports
}

func loadOneMCPServerWithRetry(
	ctx context.Context,
	index int,
	server config.MCPServerConfig,
	mcpConfig config.MCPConfig,
	loadID string,
) mcpServerLoadReport {
	report := mcpServerLoadReport{
		ServerName:        server.Name,
		ServerDescription: server.Description,
		Index:             index,
	}
	if server.Transport != "stdio" {
		report.Err = fmt.Errorf("unsupported MCP transport %q", server.Transport)
		return report
	}

	var lastErr error
	for attempt := 1; attempt <= mcpLoadMaxAttempts; attempt++ {
		report.Attempts = attempt
		client, tools, err := loadOneMCPServerAttempt(ctx, server, mcpConfig, loadID)
		if err != nil {
			lastErr = err
			continue
		}

		report.Client = client
		report.Tools = tools
		report.Success = true
		return report
	}

	report.Err = lastErr
	return report
}

func loadOneMCPServerAttempt(
	ctx context.Context,
	server config.MCPServerConfig,
	mcpConfig config.MCPConfig,
	loadID string,
) (*mcpadapter.Client, []tool.Tool, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, mcpLoadAttemptTimeout)
	defer cancel()

	env := append([]string{}, server.Env...)
	env = append(env, "UIT_HUB_MCP_LOAD_ID="+loadID)

	client, err := mcpadapter.NewStdioClient(attemptCtx, mcpadapter.StdioConfig{
		ServerName: server.Name,
		Command:    server.Command,
		Args:       server.Args,
		Env:        env,
		WorkDir:    server.WorkDir,
		Client: mcpadapter.ClientIdentity{
			Name:    mcpConfig.ClientName,
			Version: mcpConfig.ClientVersion,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	tools, err := client.Tools(attemptCtx)
	if err != nil {
		closeErr := client.Close()
		if closeErr != nil {
			return nil, nil, fmt.Errorf("list tools: %w; close client: %v", err, closeErr)
		}
		return nil, nil, err
	}
	return client, tools, nil
}

func closeMCPClient(ctx context.Context, logger *slog.Logger, client *mcpadapter.Client) {
	if client == nil {
		return
	}
	if err := client.Close(); err != nil {
		logger.DebugContext(ctx, "Close MCP client failed", "error", err)
	}
}
