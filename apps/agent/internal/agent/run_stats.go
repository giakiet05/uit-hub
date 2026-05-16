package agent

import "time"

type RunStats struct {
	StartedAt    time.Time
	FinishedAt   time.Time
	Rounds       int
	LLMCalls     int
	ToolCalls    int
	ToolFailures int
	InputTokens  int
	OutputTokens int
}

func NewRunStats() *RunStats {
	return &RunStats{
		StartedAt: time.Now(),
	}
}

func (s *RunStats) Finish() {
	s.FinishedAt = time.Now()
}

func (s RunStats) Duration() time.Duration {
	if s.FinishedAt.IsZero() {
		return time.Since(s.StartedAt)
	}
	return s.FinishedAt.Sub(s.StartedAt)
}
