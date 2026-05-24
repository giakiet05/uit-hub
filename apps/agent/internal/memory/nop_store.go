package memory

import "context"

// NopStore disables memory while preserving the Store interface.
type NopStore struct{}

// NewNopStore creates a disabled memory store.
func NewNopStore() *NopStore {
	return &NopStore{}
}

// List returns no memories when memory is disabled.
func (s *NopStore) List(ctx context.Context) ([]IndexEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []IndexEntry{}, nil
}

// Read fails because memory is disabled.
func (s *NopStore) Read(ctx context.Context, id string) (Memory, error) {
	if err := ctx.Err(); err != nil {
		return Memory{}, err
	}
	return Memory{}, ErrDisabled
}

// Write fails because memory is disabled.
func (s *NopStore) Write(ctx context.Context, input WriteInput) (Memory, error) {
	if err := ctx.Err(); err != nil {
		return Memory{}, err
	}
	return Memory{}, ErrDisabled
}
