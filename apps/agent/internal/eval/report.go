package eval

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WriteReport writes result.json and trace.md for one eval run.
func WriteReport(root string, testCase Case, result Result, trace Trace) (string, error) {
	runDir := filepath.Join(root, safePathPart(testCase.ID), time.Now().UTC().Format("20060102T150405Z"))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return "", err
	}

	resultData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(runDir, "result.json"), resultData, 0o600); err != nil {
		return "", err
	}

	traceData := []byte(renderTrace(testCase, result, trace))
	if err := os.WriteFile(filepath.Join(runDir, "trace.md"), traceData, 0o600); err != nil {
		return "", err
	}

	return runDir, nil
}

func renderTrace(testCase Case, result Result, trace Trace) string {
	var builder strings.Builder
	builder.WriteString("# Eval Trace\n\n")
	builder.WriteString(fmt.Sprintf("- Case: `%s`\n", testCase.ID))
	builder.WriteString(fmt.Sprintf("- Title: %s\n", testCase.Title))
	builder.WriteString(fmt.Sprintf("- Passed: %t\n", result.Passed))
	builder.WriteString(fmt.Sprintf("- Terminal: `%s`\n", result.TerminalReason))
	builder.WriteString(fmt.Sprintf("- Rounds: %d\n", result.RoundCount))
	builder.WriteString(fmt.Sprintf("- Tool calls: %d\n", result.ToolCalls))
	builder.WriteString(fmt.Sprintf("- LLM calls: %d\n\n", result.LLMCalls))

	builder.WriteString("## Prompt\n\n")
	builder.WriteString(testCase.Prompt)
	builder.WriteString("\n\n## Events\n\n")
	for _, event := range trace.Events {
		builder.WriteString(fmt.Sprintf("- `%s` %s\n", event.Time.Format(time.RFC3339), event.Text))
	}

	builder.WriteString("\n## Checks\n\n")
	for _, check := range result.CheckResults {
		status := "FAIL"
		if check.Passed {
			status = "PASS"
		}
		builder.WriteString(fmt.Sprintf("- %s `%s`: %s\n", status, check.Name, check.Details))
	}

	builder.WriteString("\n## Final Answer\n\n")
	builder.WriteString(result.FinalAnswer)
	builder.WriteString("\n")
	return builder.String()
}

func safePathPart(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if char >= 'a' && char <= 'z' ||
			char >= 'A' && char <= 'Z' ||
			char >= '0' && char <= '9' ||
			char == '-' ||
			char == '_' {
			builder.WriteRune(char)
			continue
		}
		builder.WriteByte('_')
	}
	return builder.String()
}
