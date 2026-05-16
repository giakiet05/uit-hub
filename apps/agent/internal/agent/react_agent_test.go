package agent

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm/echo"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestReActAgentRunsProvider(t *testing.T) {
	agent := NewReActAgent(ReActConfig{
		Provider: echo.NewProvider(),
		Prompts:  prompt.NewBuilder("system"),
		Logger:   logging.NewNopLogger(),
	})
	session := runtime.NewSession()

	message, err := agent.Run(context.Background(), session, "xin chao")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := conversation.Text(message), "received: xin chao"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
}

func TestReActAgentLogsRunStats(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	agent := NewReActAgent(ReActConfig{
		Provider: usageProvider{},
		Prompts:  prompt.NewBuilder("system"),
		Logger:   logger,
	})
	session := runtime.NewSession()

	_, err := agent.Run(context.Background(), session, "track stats")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	output := logs.String()
	expectedParts := []string{
		`msg="Agent run stats"`,
		"event=react.run.stats",
		"rounds=1",
		"llm_calls=1",
		"tool_calls=0",
		"input_tokens=12",
		"output_tokens=7",
	}
	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("stats log missing %q in %q", part, output)
		}
	}
}

func TestReActAgentHidesRawToolErrorsFromObservation(t *testing.T) {
	registry, err := tool.NewRegistry(errorTool{})
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	provider := &toolErrorProvider{}
	agent := NewReActAgent(ReActConfig{
		Provider: provider,
		Prompts:  prompt.NewBuilder("system"),
		Tools:    registry,
		Logger:   logging.NewNopLogger(),
	})
	session := runtime.NewSession()

	message, err := agent.Run(context.Background(), session, "test failing tool")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := conversation.Text(message), "handled failure"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if strings.Contains(provider.observation, "database password leaked") {
		t.Fatalf("raw tool error leaked into observation: %q", provider.observation)
	}
	if !strings.Contains(provider.observation, "tool error: execution failed") {
		t.Fatalf("sanitized tool error missing from observation: %q", provider.observation)
	}
}

func TestReActAgentTimesOutToolCalls(t *testing.T) {
	registry, err := tool.NewRegistry(blockingTool{})
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	provider := &toolTimeoutProvider{}
	agent := NewReActAgent(ReActConfig{
		Provider:    provider,
		Prompts:     prompt.NewBuilder("system"),
		Tools:       registry,
		Logger:      logging.NewNopLogger(),
		ToolTimeout: time.Millisecond,
	})
	session := runtime.NewSession()

	message, err := agent.Run(context.Background(), session, "test slow tool")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := conversation.Text(message), "handled timeout"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if !strings.Contains(provider.observation, "tool error: execution timed out") {
		t.Fatalf("timeout observation missing: %q", provider.observation)
	}
}

func TestReActAgentExecutesToolCalls(t *testing.T) {
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
	agent := NewReActAgent(ReActConfig{
		Provider: provider,
		Prompts:  prompt.NewBuilder("system"),
		Tools:    registry,
		Logger:   logging.NewNopLogger(),
	})
	session := runtime.NewSession()

	message, err := agent.Run(context.Background(), session, "test tools")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := conversation.Text(message), "all tools completed"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if got, want := provider.calls, 2; got != want {
		t.Fatalf("provider calls = %d, want %d", got, want)
	}
}

type scriptedToolProvider struct {
	calls int
}

type usageProvider struct{}

func (p usageProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("done", nil),
		Usage: llm.Usage{
			InputTokens:  12,
			OutputTokens: 7,
		},
	}, nil
}

func (p *scriptedToolProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	if p.calls == 1 {
		if got, want := len(request.Tools), 3; got != want {
			return llm.GenerateResponse{}, fmt.Errorf("tool count = %d, want %d", got, want)
		}
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
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
			}),
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
		Message: conversation.NewAssistantMessage("all tools completed", nil),
	}, nil
}

type toolErrorProvider struct {
	calls       int
	observation string
}

func (p *toolErrorProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	if p.calls == 1 {
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
				{
					ID:   "call-error",
					Name: "error_tool",
				},
			}),
		}, nil
	}

	p.observation = toolTranscript(request.Messages)
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("handled failure", nil),
	}, nil
}

type toolTimeoutProvider struct {
	calls       int
	observation string
}

func (p *toolTimeoutProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	if p.calls == 1 {
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
				{
					ID:   "call-timeout",
					Name: "blocking_tool",
				},
			}),
		}, nil
	}

	p.observation = toolTranscript(request.Messages)
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("handled timeout", nil),
	}, nil
}

type errorTool struct{}

func (t errorTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "error_tool",
		Description: "Always fails.",
		InputSchema: map[string]any{
			"type": "object",
		},
	}
}

func (t errorTool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	return tool.Result{}, errors.New("database password leaked")
}

type blockingTool struct{}

func (t blockingTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "blocking_tool",
		Description: "Blocks until context cancellation.",
		InputSchema: map[string]any{
			"type": "object",
		},
	}
}

func (t blockingTool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	<-ctx.Done()
	return tool.Result{}, ctx.Err()
}

func toolTranscript(messages []conversation.Message) string {
	var builder strings.Builder
	for _, message := range messages {
		toolResult, ok := message.(conversation.ToolResultMessage)
		if !ok {
			continue
		}
		builder.WriteString(toolResult.ToolCallID)
		builder.WriteString("=")
		builder.WriteString(conversation.Text(toolResult))
		builder.WriteString("\n")
	}
	return builder.String()
}
