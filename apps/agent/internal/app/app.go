package app

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/cli"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/echo"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/openai"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	err := config.LoadEnv()
	logger := logging.NewLogger(stderr)
	if err != nil {
		logger.DebugContext(ctx, "No .env file loaded", "error", err)
	} else {
		logger.DebugContext(ctx, ".env file loaded")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.DebugContext(ctx, "Agent config load failed", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent config loaded", "llm_provider", cfg.Provider, "openai_model", cfg.OpenAI.Model, "openai_base_url", cfg.OpenAI.BaseURL)
	logger.DebugContext(ctx, "Starting agent app", "llm_provider", cfg.Provider)

	provider, err := newProvider(cfg, logger)
	if err != nil {
		return err
	}

	session := runtime.NewSession()
	prompts := prompt.NewBuilder(cfg.SystemPrompt)
	tools, err := newToolRegistry()
	if err != nil {
		logger.DebugContext(ctx, "Tool registry initialization failed", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Tool registry initialized", "tool_count", len(tools.Definitions()))

	loop := agent.NewLoop(provider, prompts, logger, tools)
	runner := cli.NewRunner(session, loop, stdin, stdout, stderr, logger)

	if err := runner.Run(ctx, args); err != nil {
		logger.DebugContext(ctx, "Agent app stopped with error", "error", err)
		return err
	}
	logger.DebugContext(ctx, "Agent app stopped")
	return nil
}

func newToolRegistry() (*tool.Registry, error) {
	return tool.NewRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
		localtool.NewCurrentTime(),
	)
}

func newProvider(cfg config.Config, logger *slog.Logger) (llm.Provider, error) {
	switch cfg.Provider {
	case "openai":
		logger.Debug("Initializing OpenAI LLM provider", "model", cfg.OpenAI.Model, "base_url", cfg.OpenAI.BaseURL)
		return openai.NewProvider(openai.Config{
			APIKey:  cfg.OpenAI.APIKey,
			Model:   cfg.OpenAI.Model,
			BaseURL: cfg.OpenAI.BaseURL,
		}, logger), nil
	case "echo":
		logger.Debug("Initializing echo LLM provider")
		return echo.NewProvider(), nil
	default:
		logger.Debug("Unsupported LLM provider requested", "llm_provider", cfg.Provider)
		return nil, errors.New("unsupported LLM_PROVIDER: " + cfg.Provider)
	}
}
