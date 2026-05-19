package localtool

import (
	"context"
	"fmt"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// CurrentTime returns the current time using an injectable clock.
type CurrentTime struct {
	now func() time.Time
}

// NewCurrentTime creates a time tool backed by time.Now.
func NewCurrentTime() *CurrentTime {
	return &CurrentTime{now: time.Now}
}

// NewCurrentTimeWithClock creates a time tool with a deterministic clock for
// tests.
func NewCurrentTimeWithClock(now func() time.Time) *CurrentTime {
	return &CurrentTime{now: now}
}

// Definition describes the current-time tool schema.
func (t *CurrentTime) Definition() tool.Definition {
	return tool.Definition{
		Name:        "current_time",
		Description: "Return the current time in RFC3339 format for a timezone.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"timezone": map[string]any{
					"type":        "string",
					"description": "IANA timezone name, for example Asia/Ho_Chi_Minh. Defaults to UTC.",
				},
			},
			"required": []string{"timezone"},
		},
	}
}

// Execute returns the current time in the requested IANA timezone.
func (t *CurrentTime) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	timezone, err := stringArg(call.Arguments, "timezone")
	if err != nil {
		return tool.Result{}, err
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return tool.Result{}, fmt.Errorf("load timezone: %w", err)
	}

	now := time.Now
	if t.now != nil {
		now = t.now
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: now().In(location).Format(time.RFC3339),
	}, nil
}
