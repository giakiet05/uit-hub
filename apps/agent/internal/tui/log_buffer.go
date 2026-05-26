package tui

import (
	"strings"
	"sync"
)

// LogBuffer is an in-memory io.Writer that keeps the most recent log lines for
// the TUI.
type LogBuffer struct {
	mu    sync.Mutex
	lines []string
	limit int
}

// NewLogBuffer creates a bounded log buffer.
func NewLogBuffer(limit int) *LogBuffer {
	if limit <= 0 {
		limit = 500
	}
	return &LogBuffer{
		lines: []string{},
		limit: limit,
	}
}

// Write appends log bytes to the buffer and trims older lines beyond the limit.
func (b *LogBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	chunks := strings.Split(strings.TrimRight(string(p), "\n"), "\n")
	for _, chunk := range chunks {
		if chunk == "" {
			continue
		}
		b.lines = append(b.lines, chunk)
	}
	if overflow := len(b.lines) - b.limit; overflow > 0 {
		b.lines = append([]string{}, b.lines[overflow:]...)
	}
	return len(p), nil
}

// Lines returns a defensive copy of buffered log lines.
func (b *LogBuffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := make([]string, len(b.lines))
	copy(lines, b.lines)
	return lines
}
