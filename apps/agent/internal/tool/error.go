package tool

import (
	"context"
	"errors"
	"fmt"
)

// ErrorType classifies tool execution failures for sanitized model
// observations and future metrics.
type ErrorType string

const (
	ErrorTypeNotFound        ErrorType = "not_found"
	ErrorTypeTimeout         ErrorType = "timeout"
	ErrorTypeExecutionFailed ErrorType = "execution_failed"
	ErrorTypeValidation      ErrorType = "validation_failed"
)

// ToolError wraps an internal tool failure with a stable public category.
type ToolError struct {
	Type ErrorType
	Name string
	Err  error
}

// Error returns a compact developer-facing error string.
func (e ToolError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("tool %q failed: %s", e.Name, e.Type)
	}
	return fmt.Sprintf("tool %q failed: %s: %v", e.Name, e.Type, e.Err)
}

// Unwrap exposes the wrapped internal error for errors.Is/errors.As.
func (e ToolError) Unwrap() error {
	return e.Err
}

// ClassifyError converts an arbitrary tool error into a ToolError.
func ClassifyError(name string, err error) ToolError {
	var existing ToolError
	if errors.As(err, &existing) {
		// Preserve the existing ToolError but optionally update the name if it's empty
		if existing.Name == "" {
			existing.Name = name
		}
		return existing
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return ToolError{Type: ErrorTypeNotFound, Name: name, Err: err}
	case errors.Is(err, context.DeadlineExceeded):
		return ToolError{Type: ErrorTypeTimeout, Name: name, Err: err}
	default:
		return ToolError{Type: ErrorTypeExecutionFailed, Name: name, Err: err}
	}
}

// ErrorObservation converts internal tool failures into sanitized observations
// that can be safely shown back to the model.
func ErrorObservation(err error) string {
	var toolErr ToolError
	if !errors.As(err, &toolErr) {
		toolErr = ClassifyError("", err)
	}

	switch toolErr.Type {
	case ErrorTypeTimeout:
		return "tool error: execution timed out"
	case ErrorTypeNotFound:
		return "tool error: requested tool is not available"
	case ErrorTypeValidation:
		return fmt.Sprintf("tool error: input validation failed: %v", toolErr.Err)
	default:
		return "tool error: execution failed"
	}
}
