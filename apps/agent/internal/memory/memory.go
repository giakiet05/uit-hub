// Package memory defines long-term agent memory independent of storage backend.
package memory

import (
	"context"
	"errors"
	"time"
)

// Type classifies long-term memories so the agent does not save random facts.
type Type string

const (
	TypeUser      Type = "user"
	TypeFeedback  Type = "feedback"
	TypeProject   Type = "project"
	TypeReference Type = "reference"
)

// ErrNotFound marks a request for a memory ID that does not exist.
var ErrNotFound = errors.New("memory not found")

// ErrDisabled marks an operation against a disabled memory store.
var ErrDisabled = errors.New("memory disabled")

// Memory is one durable observation the agent may use across sessions.
type Memory struct {
	ID          string
	Type        Type
	Name        string
	Description string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IndexEntry is the lightweight table-of-contents row loaded into the prompt.
type IndexEntry struct {
	ID          string
	Type        Type
	Name        string
	Description string
	UpdatedAt   time.Time
}

// WriteInput creates or updates a memory. Empty ID means create a new memory.
type WriteInput struct {
	ID          string
	Type        Type
	Name        string
	Description string
	Content     string
}

// Store hides the storage backend. FileStore is for CLI; DBStore can replace it
// later.
type Store interface {
	List(ctx context.Context) ([]IndexEntry, error)
	Read(ctx context.Context, id string) (Memory, error)
	Write(ctx context.Context, input WriteInput) (Memory, error)
}
