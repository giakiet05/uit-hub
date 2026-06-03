package config

import (
	"log/slog"
	"os"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/openai"
	"gopkg.in/yaml.v3"
)

// YAMLModel represents a single model entry in the YAML configuration.
type YAMLModel struct {
	Name             string `yaml:"name"`
	MaxContextTokens int    `yaml:"max_context_tokens"`
}

// YAMLProvider represents a provider entry in the YAML configuration.
type YAMLProvider struct {
	APIKeyEnv      string      `yaml:"api_key_env"`
	BaseURLEnv     string      `yaml:"base_url_env"`
	DefaultBaseURL string      `yaml:"default_base_url"`
	Models         []YAMLModel `yaml:"models"`
}

// YAMLModelsConfig is the root structure for models.yaml.
type YAMLModelsConfig struct {
	Providers map[string]YAMLProvider `yaml:"providers"`
}

// LoadLLMRegistry reads the YAML configuration, instantiates all providers and models,
// and registers them into an llm.Registry.
func LoadLLMRegistry(configPath string, logger *slog.Logger) (*llm.Registry, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config YAMLModelsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	registry := llm.NewRegistry()

	for providerName, yamlProv := range config.Providers {
		apiKey := os.Getenv(yamlProv.APIKeyEnv)
		baseURL := os.Getenv(yamlProv.BaseURLEnv)
		if baseURL == "" {
			baseURL = yamlProv.DefaultBaseURL
		}

		var provider llm.Provider
		switch providerName {
		case "openai":
			provider = openai.NewProvider(apiKey, baseURL, logger)
		// Add other providers here later (e.g., gemini, anthropic)
		default:
			logger.Warn("Unknown provider in models.yaml", "provider", providerName)
			continue
		}

		for _, yamlMod := range yamlProv.Models {
			model, err := provider.CreateModel(yamlMod.Name, yamlMod.MaxContextTokens)
			if err != nil {
				logger.Warn("Failed to create model", "provider", providerName, "model", yamlMod.Name, "error", err)
				continue
			}
			registry.RegisterModel(providerName, model)
			logger.Debug("Registered model", "provider", providerName, "model", yamlMod.Name)
		}
	}

	return registry, nil
}
