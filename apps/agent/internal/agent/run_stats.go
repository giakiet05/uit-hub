package agent

import "time"

// RunStats records the high-level counters for a single agent run.
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

// NewRunStats starts a stats record with the current time.
func NewRunStats() *RunStats {
	return &RunStats{
		StartedAt: time.Now(),
	}
}

// Finish marks the run as completed.
func (s *RunStats) Finish() {
	s.FinishedAt = time.Now()
}

// Duration returns elapsed time for a completed run, or live elapsed time while
// the run is still active.
func (s RunStats) Duration() time.Duration {
	if s.FinishedAt.IsZero() {
		return time.Since(s.StartedAt)
	}
	return s.FinishedAt.Sub(s.StartedAt)
}
