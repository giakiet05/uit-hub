package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/echo"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestLoopRunsProvider(t *testing.T) {
	loop := NewLoop(echo.NewProvider(), prompt.NewBuilder("system"), logging.NewNopLogger())
	session := runtime.NewSession()

	message, err := loop.Run(context.Background(), session, "xin chao")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := message.ContentText(), "received: xin chao"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
}

func TestLoopExecutesToolCalls(t *testing.T) {
	registry, err := tool.NewRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
		localtool.NewCurrentTimeWithClock(func() time.Time {
			return time.Date(2026, 5, 14, 10, 30, 0, 0, time.UTC)
		}),
	)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	provider := &scriptedToolProvider{}
	loop := NewLoop(provider, prompt.NewBuilder("system"), logging.NewNopLogger(), registry)
	session := runtime.NewSession()

	message, err := loop.Run(context.Background(), session, "test tools")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := message.ContentText(), "all tools completed"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if got, want := provider.calls, 2; got != want {
		t.Fatalf("provider calls = %d, want %d", got, want)
	}
}

type scriptedToolProvider struct {
	calls int
}

func (p *scriptedToolProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	if p.calls == 1 {
		if got, want := len(request.Tools), 3; got != want {
			return llm.GenerateResponse{}, fmt.Errorf("tool count = %d, want %d", got, want)
		}
		return llm.GenerateResponse{
			Message: llm.Message{
				Role: llm.RoleAssistant,
				ToolCalls: []llm.ToolCall{
					{
						ID:   "call-echo",
						Name: "echo",
						Arguments: map[string]any{
							"text": "hello tool",
						},
					},
					{
						ID:   "call-calc",
						Name: "calculator",
						Arguments: map[string]any{
							"operation": "multiply",
							"a":         6,
							"b":         7,
						},
					},
					{
						ID:   "call-time",
						Name: "current_time",
						Arguments: map[string]any{
							"timezone": "Asia/Ho_Chi_Minh",
						},
					},
				},
			},
			FinishReason: llm.FinishReasonToolCall,
		}, nil
	}

	transcript := toolTranscript(request.Messages)
	if !strings.Contains(transcript, "call-echo=hello tool") {
		return llm.GenerateResponse{}, fmt.Errorf("missing echo result in %q", transcript)
	}
	if !strings.Contains(transcript, "call-calc=42") {
		return llm.GenerateResponse{}, fmt.Errorf("missing calculator result in %q", transcript)
	}
	if !strings.Contains(transcript, "call-time=2026-05-14T17:30:00+07:00") {
		return llm.GenerateResponse{}, fmt.Errorf("missing current_time result in %q", transcript)
	}

	return llm.GenerateResponse{
		Message:      llm.NewTextMessage(llm.RoleAssistant, "all tools completed"),
		FinishReason: llm.FinishReasonStop,
	}, nil
}

func toolTranscript(messages []llm.Message) string {
	var builder strings.Builder
	for _, message := range messages {
		if message.Role != llm.RoleTool {
			continue
		}
		builder.WriteString(message.ToolCallID)
		builder.WriteString("=")
		builder.WriteString(message.ContentText())
		builder.WriteString("\n")
	}
	return builder.String()
}
