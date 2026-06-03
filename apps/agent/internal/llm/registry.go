package llm

import (
	"errors"
	"fmt"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrModelNotFound    = errors.New("model not found")
)

// Registry acts as a central repository for instantiated language models.
type Registry struct {
	// models maps provider name -> model name -> Model instance.
	models map[string]map[string]Model
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		models: make(map[string]map[string]Model),
	}
}

// RegisterModel adds a fully instantiated model to the registry.
func (r *Registry) RegisterModel(providerName string, model Model) {
	if _, ok := r.models[providerName]; !ok {
		r.models[providerName] = make(map[string]Model)
	}
	r.models[providerName][model.Name()] = model
}

// GetModel retrieves a model by its provider and name.
func (r *Registry) GetModel(providerName, modelName string) (Model, error) {
	providerModels, ok := r.models[providerName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, providerName)
	}

	model, ok := providerModels[modelName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrModelNotFound, modelName)
	}

	return model, nil
}
