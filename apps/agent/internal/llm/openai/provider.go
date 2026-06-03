// Package openai adapts the OpenAI Responses API to the internal LLM provider
// interface.
package openai

import (
	"log/slog"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	openaisdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Provider implements llm.Provider as a factory for OpenAI models.
type Provider struct {
	llm.BaseProvider
	client openaisdk.Client
	logger *slog.Logger
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return string(llm.ProviderTypeOpenAI)
}

// CreateModel instantiates a new OpenAI model.
func (p *Provider) CreateModel(modelName string, maxContext int) (llm.Model, error) {
	return &Model{
		BaseModel: llm.NewBaseModel(modelName, maxContext, p.logger),
		client:    p.client,
	}, nil
}

// NewProvider creates an OpenAI provider backed by the official OpenAI Go SDK.
func NewProvider(apiKey, baseURL string, logger *slog.Logger) *Provider {
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	options := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		options = append(options, option.WithBaseURL(strings.TrimRight(baseURL, "/")))
	}

	return &Provider{
		BaseProvider: llm.NewBaseProvider(apiKey, baseURL),
		client:       openaisdk.NewClient(options...),
		logger:       logger,
	}
}


