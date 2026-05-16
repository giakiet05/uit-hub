package tui

import (
	"strings"
	"sync"
)

type LogBuffer struct {
	mu    sync.Mutex
	lines []string
	limit int
}

func NewLogBuffer(limit int) *LogBuffer {
	if limit <= 0 {
		limit = 500
	}
	return &LogBuffer{
		lines: []string{},
		limit: limit,
	}
}

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

func (b *LogBuffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := make([]string, len(b.lines))
	copy(lines, b.lines)
	return lines
}
