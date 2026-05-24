package memory

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileStoreWritesListsAndReadsMemory(t *testing.T) {
	store := NewFileStore(t.TempDir())
	store.clock = fixedClock(time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC))

	saved, err := store.Write(context.Background(), WriteInput{
		Type:        TypeFeedback,
		Name:        "Response Style",
		Description: "User prefers concise Vietnamese answers.",
		Content:     "Answer directly and mention tradeoffs.",
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if saved.ID != "feedback_response_style" {
		t.Fatalf("saved.ID = %q, want generated slug", saved.ID)
	}

	entries, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got, want := len(entries), 1; got != want {
		t.Fatalf("len(entries) = %d, want %d", got, want)
	}
	if entries[0].ID != saved.ID {
		t.Fatalf("entries[0].ID = %q, want %q", entries[0].ID, saved.ID)
	}

	read, err := store.Read(context.Background(), saved.ID)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if read.Content != "Answer directly and mention tradeoffs." {
		t.Fatalf("read.Content = %q", read.Content)
	}

	index, err := os.ReadFile(filepath.Join(store.root(), memoryIndexFile))
	if err != nil {
		t.Fatalf("ReadFile(MEMORY.md) error = %v", err)
	}
	if !strings.Contains(string(index), "Response Style") {
		t.Fatalf("MEMORY.md missing entry: %q", string(index))
	}
}

func TestFileStoreUpdatePreservesCreatedAt(t *testing.T) {
	store := NewFileStore(t.TempDir())
	createdAt := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 5, 21, 11, 0, 0, 0, time.UTC)
	store.clock = fixedClock(createdAt)

	saved, err := store.Write(context.Background(), WriteInput{
		Type:        TypeUser,
		Name:        "Role",
		Description: "User is building an agent.",
		Content:     "User is building a Go agent.",
	})
	if err != nil {
		t.Fatalf("Write() create error = %v", err)
	}

	store.clock = fixedClock(updatedAt)
	updated, err := store.Write(context.Background(), WriteInput{
		ID:          saved.ID,
		Type:        TypeUser,
		Name:        "Role",
		Description: "User is building a reusable agent.",
		Content:     "User is building a reusable Go agent.",
	})
	if err != nil {
		t.Fatalf("Write() update error = %v", err)
	}
	if !updated.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %s, want %s", updated.CreatedAt, createdAt)
	}
	if !updated.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want %s", updated.UpdatedAt, updatedAt)
	}
}

func TestFileStoreWriteWithEmptyIDUpdatesGeneratedSlug(t *testing.T) {
	store := NewFileStore(t.TempDir())
	createdAt := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 5, 21, 11, 0, 0, 0, time.UTC)
	store.clock = fixedClock(createdAt)

	saved, err := store.Write(context.Background(), WriteInput{
		Type:        TypeUser,
		Name:        "User profile and preferences",
		Description: "Initial user profile.",
		Content:     "Name: Kiet",
	})
	if err != nil {
		t.Fatalf("Write() create error = %v", err)
	}

	store.clock = fixedClock(updatedAt)
	updated, err := store.Write(context.Background(), WriteInput{
		Type:        TypeUser,
		Name:        "User profile and preferences",
		Description: "Updated user profile.",
		Content:     "Name: Kiet\nMajor: Software Engineering",
	})
	if err != nil {
		t.Fatalf("Write() generated-slug update error = %v", err)
	}
	if updated.ID != saved.ID {
		t.Fatalf("updated.ID = %q, want %q", updated.ID, saved.ID)
	}
	if !updated.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %s, want %s", updated.CreatedAt, createdAt)
	}
	if !updated.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want %s", updated.UpdatedAt, updatedAt)
	}
	if !strings.Contains(updated.Content, "Software Engineering") {
		t.Fatalf("updated.Content = %q", updated.Content)
	}
}

func TestFileStoreRejectsInvalidID(t *testing.T) {
	store := NewFileStore(t.TempDir())

	_, err := store.Read(context.Background(), "../escape")
	if err == nil {
		t.Fatal("Read() error = nil, want invalid id error")
	}
}

func TestFileStoreReadMissingMemory(t *testing.T) {
	store := NewFileStore(t.TempDir())

	_, err := store.Read(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Read() error = %v, want ErrNotFound", err)
	}
}

func TestFormatIndexEmpty(t *testing.T) {
	if got, want := FormatIndex(nil), "No saved memories."; got != want {
		t.Fatalf("FormatIndex(nil) = %q, want %q", got, want)
	}
}

func fixedClock(now time.Time) func() time.Time {
	return func() time.Time {
		return now
	}
}
