package memory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const memoryIndexFile = "MEMORY.md"

var validID = regexp.MustCompile(`^[a-z0-9_-]+$`)

// FileStore stores memories as Markdown files for local CLI development.
type FileStore struct {
	rootDir string
	clock   func() time.Time
}

// NewFileStore creates a Markdown-backed memory store.
func NewFileStore(rootDir string) *FileStore {
	return &FileStore{
		rootDir: rootDir,
		clock:   time.Now,
	}
}

// List scans memory files and returns only frontmatter metadata.
func (s *FileStore) List(ctx context.Context) ([]IndexEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := s.ensureRoot(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.root())
	if err != nil {
		return nil, fmt.Errorf("read memory directory: %w", err)
	}

	index := make([]IndexEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || entry.Name() == memoryIndexFile {
			continue
		}
		memory, err := s.readFile(filepath.Join(s.root(), entry.Name()))
		if err != nil {
			return nil, err
		}
		index = append(index, IndexEntry{
			ID:          memory.ID,
			Type:        memory.Type,
			Name:        memory.Name,
			Description: memory.Description,
			UpdatedAt:   memory.UpdatedAt,
		})
	}

	slices.SortFunc(index, func(left IndexEntry, right IndexEntry) int {
		return strings.Compare(left.ID, right.ID)
	})
	return index, nil
}

// Read loads one full memory by ID.
func (s *FileStore) Read(ctx context.Context, id string) (Memory, error) {
	if err := ctx.Err(); err != nil {
		return Memory{}, err
	}

	path, err := s.pathForID(id)
	if err != nil {
		return Memory{}, err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return Memory{}, ErrNotFound
	}
	if err != nil {
		return Memory{}, fmt.Errorf("stat memory file: %w", err)
	}
	return s.readFile(path)
}

// Write creates or updates one memory file and regenerates MEMORY.md.
func (s *FileStore) Write(ctx context.Context, input WriteInput) (Memory, error) {
	if err := ctx.Err(); err != nil {
		return Memory{}, err
	}
	if err := validateWriteInput(input); err != nil {
		return Memory{}, err
	}
	if err := s.ensureRoot(); err != nil {
		return Memory{}, err
	}

	now := s.now()
	id := strings.TrimSpace(input.ID)
	if id == "" {
		id = slugID(input.Type, input.Name)
	}
	path, err := s.pathForID(id)
	if err != nil {
		return Memory{}, err
	}

	createdAt := now
	if input.ID != "" {
		existing, err := s.Read(ctx, id)
		if err != nil {
			return Memory{}, err
		}
		createdAt = existing.CreatedAt
	} else if existing, err := s.Read(ctx, id); err == nil {
		createdAt = existing.CreatedAt
	} else if !errors.Is(err, ErrNotFound) {
		return Memory{}, err
	}

	memory := Memory{
		ID:          id,
		Type:        input.Type,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Content:     strings.TrimSpace(input.Content),
		CreatedAt:   createdAt,
		UpdatedAt:   now,
	}
	if err := os.WriteFile(path, []byte(renderMemoryFile(memory)), 0o644); err != nil {
		return Memory{}, fmt.Errorf("write memory file: %w", err)
	}
	if err := s.writeIndex(ctx); err != nil {
		return Memory{}, err
	}

	return memory, nil
}

// root returns the configured memory directory or the default local directory.
func (s *FileStore) root() string {
	if strings.TrimSpace(s.rootDir) == "" {
		return "memory"
	}
	return s.rootDir
}

// now returns the current time through an injectable clock so timestamp tests
// are deterministic.
func (s *FileStore) now() time.Time {
	if s.clock == nil {
		return time.Now().UTC()
	}
	return s.clock().UTC()
}

// ensureRoot creates the memory directory before any read/write operation that
// expects it to exist.
func (s *FileStore) ensureRoot() error {
	if err := os.MkdirAll(s.root(), 0o755); err != nil {
		return fmt.Errorf("create memory directory: %w", err)
	}
	return nil
}

// pathForID converts a memory ID into an absolute file path under the memory
// root and rejects IDs that could escape the root.
func (s *FileStore) pathForID(id string) (string, error) {
	id = strings.TrimSpace(id)
	if !validID.MatchString(id) {
		return "", fmt.Errorf("invalid memory id %q", id)
	}

	root, err := filepath.Abs(s.root())
	if err != nil {
		return "", fmt.Errorf("resolve memory root: %w", err)
	}
	path := filepath.Join(root, id+".md")
	if filepath.Dir(path) != root {
		return "", errors.New("memory path escapes root")
	}
	return path, nil
}

// readFile reads and parses one Markdown memory file.
func (s *FileStore) readFile(path string) (Memory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Memory{}, fmt.Errorf("read memory file: %w", err)
	}
	memory, err := parseMemoryFile(string(data))
	if err != nil {
		return Memory{}, fmt.Errorf("parse memory file %q: %w", filepath.Base(path), err)
	}
	return memory, nil
}

// writeIndex regenerates MEMORY.md from the current memory metadata.
func (s *FileStore) writeIndex(ctx context.Context) error {
	entries, err := s.List(ctx)
	if err != nil {
		return err
	}

	path := filepath.Join(s.root(), memoryIndexFile)
	if err := os.WriteFile(path, []byte(FormatIndex(entries)+"\n"), 0o644); err != nil {
		return fmt.Errorf("write memory index: %w", err)
	}
	return nil
}

// validateWriteInput checks the tool-facing memory write payload before it is
// converted into a file.
func validateWriteInput(input WriteInput) error {
	if input.ID != "" && !validID.MatchString(input.ID) {
		return fmt.Errorf("invalid memory id %q", input.ID)
	}
	if !validType(input.Type) {
		return fmt.Errorf("invalid memory type %q", input.Type)
	}
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("memory name is required")
	}
	if strings.TrimSpace(input.Description) == "" {
		return errors.New("memory description is required")
	}
	if strings.TrimSpace(input.Content) == "" {
		return errors.New("memory content is required")
	}
	return nil
}

// validType reports whether a memory type belongs to the supported taxonomy.
func validType(memoryType Type) bool {
	switch memoryType {
	case TypeUser, TypeFeedback, TypeProject, TypeReference:
		return true
	default:
		return false
	}
}

// slugID creates a stable file-safe ID from memory type and display name.
func slugID(memoryType Type, name string) string {
	text := strings.ToLower(strings.TrimSpace(string(memoryType) + "_" + name))
	var builder strings.Builder
	previousSeparator := false
	for _, r := range text {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			previousSeparator = false
		case !previousSeparator:
			builder.WriteByte('_')
			previousSeparator = true
		}
	}
	return strings.Trim(builder.String(), "_")
}

type memoryFrontmatter struct {
	ID          string    `yaml:"id"`
	Type        Type      `yaml:"type"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
}

// renderMemoryFile serializes a memory as Markdown with YAML frontmatter.
func renderMemoryFile(memory Memory) string {
	header, err := yaml.Marshal(memoryFrontmatter{
		ID:          memory.ID,
		Type:        memory.Type,
		Name:        memory.Name,
		Description: memory.Description,
		CreatedAt:   memory.CreatedAt,
		UpdatedAt:   memory.UpdatedAt,
	})
	if err != nil {
		header = []byte{}
	}
	return "---\n" + string(header) + "---\n\n" + strings.TrimSpace(memory.Content) + "\n"
}

// parseMemoryFile parses Markdown with YAML frontmatter back into a Memory.
func parseMemoryFile(text string) (Memory, error) {
	if !strings.HasPrefix(text, "---\n") {
		return Memory{}, errors.New("missing frontmatter")
	}

	parts := strings.SplitN(strings.TrimPrefix(text, "---\n"), "\n---", 2)
	if len(parts) != 2 {
		return Memory{}, errors.New("unterminated frontmatter")
	}

	var header memoryFrontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &header); err != nil {
		return Memory{}, err
	}
	memory := Memory{
		ID:          strings.TrimSpace(header.ID),
		Type:        header.Type,
		Name:        strings.TrimSpace(header.Name),
		Description: strings.TrimSpace(header.Description),
		Content:     strings.TrimSpace(parts[1]),
		CreatedAt:   header.CreatedAt,
		UpdatedAt:   header.UpdatedAt,
	}
	if err := validateMemory(memory); err != nil {
		return Memory{}, err
	}
	return memory, nil
}

// validateMemory checks metadata read from disk before exposing it to callers.
func validateMemory(memory Memory) error {
	if !validID.MatchString(memory.ID) {
		return fmt.Errorf("invalid memory id %q", memory.ID)
	}
	if !validType(memory.Type) {
		return fmt.Errorf("invalid memory type %q", memory.Type)
	}
	if memory.Name == "" {
		return errors.New("memory name is required")
	}
	if memory.Description == "" {
		return errors.New("memory description is required")
	}
	return nil
}
