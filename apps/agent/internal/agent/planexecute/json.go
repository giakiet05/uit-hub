package planexecute

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func decodeStepResult(stepID string, text string, toolCalls []string, observations []string) (stepResult, error) {
	var result stepResult
	if err := decodeJSONText(text, &result); err != nil {
		return stepResult{}, fmt.Errorf("decode step result: %w", err)
	}
	if result.StepID == "" {
		result.StepID = stepID
	}
	if result.Status == "" {
		result.Status = "completed"
	}
	if len(result.ToolCalls) == 0 {
		result.ToolCalls = toolCalls
	}
	if len(result.Observations) == 0 {
		result.Observations = observations
	}
	if result.StateUpdates == nil {
		result.StateUpdates = map[string]string{}
	}
	if strings.TrimSpace(result.Summary) == "" {
		return stepResult{}, fmt.Errorf("step %q returned empty summary", stepID)
	}
	return result, nil
}

func decodeJSONText(text string, target any) error {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	if err := json.Unmarshal([]byte(text), target); err == nil {
		return nil
	}

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start < 0 || end < start {
		return errors.New("response does not contain a JSON object")
	}
	return json.Unmarshal([]byte(text[start:end+1]), target)
}

func validatePlan(plan executionPlan) error {
	if strings.TrimSpace(plan.Goal) == "" {
		return errors.New("plan goal is empty")
	}
	if len(plan.Steps) == 0 {
		return errors.New("plan has no steps")
	}
	for index, step := range plan.Steps {
		if strings.TrimSpace(step.ID) == "" {
			return fmt.Errorf("plan step %d has empty id", index+1)
		}
		if strings.TrimSpace(step.Description) == "" {
			return fmt.Errorf("plan step %q has empty description", step.ID)
		}
	}
	return nil
}

func planSummary(plan executionPlan) string {
	return fmt.Sprintf("%s (%d steps)", plan.Goal, len(plan.Steps))
}

func mustJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(data)
}
