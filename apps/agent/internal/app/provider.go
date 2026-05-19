package app

import (
	"errors"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/openai"
)

// newProvider selects the configured LLM provider.
func newProvider(cfg config.Config, logger *slog.Logger) (llm.Provider, error) {
	switch cfg.Provider {
	case llm.ProviderTypeOpenAI:
		logger.Debug("Initializing OpenAI LLM provider", "model", cfg.OpenAI.Model, "base_url", cfg.OpenAI.BaseURL)
		return openai.NewProvider(openai.Config{
			APIKey:  cfg.OpenAI.APIKey,
			Model:   cfg.OpenAI.Model,
			BaseURL: cfg.OpenAI.BaseURL,
		}, logger), nil
	default:
		logger.Debug("Unsupported LLM provider requested", "llm_provider", cfg.Provider)
		return nil, errors.New("unsupported LLM_PROVIDER: " + string(cfg.Provider))
	}
}
