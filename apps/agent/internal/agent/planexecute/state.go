package planexecute

import (
	"fmt"
	"strings"
)

func newExecutionState() executionState {
	return executionState{
		Facts: map[string]string{},
	}
}

func (s executionState) apply(result stepResult) {
	if s.Facts == nil {
		s.Facts = map[string]string{}
	}
	for key, value := range result.StateUpdates {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		s.Facts[key] = value
	}
	if strings.TrimSpace(result.Summary) != "" {
		s.Facts[fmt.Sprintf("step.%s.summary", result.StepID)] = result.Summary
	}
}

func compactResults(results []stepResult) []compactStepResult {
	compacted := make([]compactStepResult, 0, len(results))
	for _, result := range results {
		compacted = append(compacted, compactStepResult{
			StepID:       result.StepID,
			Status:       result.Status,
			Summary:      result.Summary,
			ToolCalls:    result.ToolCalls,
			StateUpdates: result.StateUpdates,
			Error:        result.Error,
		})
	}
	return compacted
}
