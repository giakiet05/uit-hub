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
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestReActAgentRunsProvider(t *testing.T) {
	agent := NewReActAgent(ReActConfig{
		Provider: echoTestProvider{},
		Prompts:  newTestPromptBuilder(),
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

func TestReActAgentSendsSystemPromptWithoutStoringItInConversation(t *testing.T) {
	provider := &captureMessagesProvider{}
	agent := NewReActAgent(ReActConfig{
		Provider: provider,
		Prompts:  newTestPromptBuilder(),
		Logger:   logging.NewNopLogger(),
	})
	session := runtime.NewSession()

	_, err := agent.Run(context.Background(), session, "xin chao")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(provider.messages) == 0 {
		t.Fatal("provider received no messages")
	}
	if _, ok := provider.messages[0].(conversation.SystemMessage); !ok {
		t.Fatalf("first provider message type = %T, want conversation.SystemMessage", provider.messages[0])
	}
	for _, message := range session.Conversation.Messages() {
		if _, ok := message.(conversation.SystemMessage); ok {
			t.Fatalf("session conversation stored system message: %#v", message)
		}
	}
}

func TestReActAgentLogsRunStats(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	agent := NewReActAgent(ReActConfig{
		Provider: usageProvider{},
		Prompts:  newTestPromptBuilder(),
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

func TestReActAgentLogsTimeline(t *testing.T) {
	registry, err := tool.NewRegistry(localtool.NewEcho())
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	agent := NewReActAgent(ReActConfig{
		Provider: &timelineProvider{},
		Prompts:  newTestPromptBuilder(),
		Tools:    registry,
		Logger:   logger,
	})
	session := runtime.NewSession()

	_, err = agent.Run(context.Background(), session, "track timeline")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	output := logs.String()
	expectedParts := []string{
		`msg="Round 1: Think"`,
		`msg="Round 1: Act"`,
		"decision_summary=",
		"Need echo to mirror the requested text",
		`msg="Round 1: Tool echo started"`,
		"arguments_preview=",
		`msg="Round 1: Observe echo"`,
		"result_preview=hello",
		`msg="Round 2: Think"`,
		`msg="Round 2: Answer"`,
	}
	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("timeline log missing %q in %q", part, output)
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
		Prompts:  newTestPromptBuilder(),
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
		Prompts:     newTestPromptBuilder(),
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
		Prompts:  newTestPromptBuilder(),
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

func TestReActAgentRunsMultipleToolRounds(t *testing.T) {
	registry, err := tool.NewRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
	)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	provider := &multiRoundToolProvider{}
	agent := NewReActAgent(ReActConfig{
		Provider: provider,
		Prompts:  newTestPromptBuilder(),
		Tools:    registry,
		Logger:   logging.NewNopLogger(),
	})
	session := runtime.NewSession()

	message, err := agent.Run(context.Background(), session, "multi round")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := conversation.Text(message), "multi-round complete"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if got, want := provider.calls, 3; got != want {
		t.Fatalf("provider calls = %d, want %d", got, want)
	}
	if got, want := countToolResults(session.Conversation.Messages()), 2; got != want {
		t.Fatalf("tool result messages = %d, want %d", got, want)
	}
}

type scriptedToolProvider struct {
	calls int
}

type usageProvider struct{}

type echoTestProvider struct{}

func (p echoTestProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	for i := len(request.Messages) - 1; i >= 0; i-- {
		if _, ok := request.Messages[i].(conversation.UserMessage); ok {
			return llm.GenerateResponse{
				Message: conversation.NewAssistantMessage("received: "+conversation.Text(request.Messages[i]), nil),
			}, nil
		}
	}
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("received: ", nil),
	}, nil
}

type captureMessagesProvider struct {
	messages []conversation.Message
}

func (p *captureMessagesProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.messages = request.Messages
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("done", nil),
	}, nil
}

func (p usageProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("done", nil),
		Usage: llm.Usage{
			InputTokens:  12,
			OutputTokens: 7,
		},
	}, nil
}

type timelineProvider struct {
	calls int
}

func (p *timelineProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	if p.calls == 1 {
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("Need echo to mirror the requested text.", []conversation.ToolCall{
				{
					ID:   "call-echo",
					Name: "echo",
					Arguments: map[string]any{
						"text": "hello",
					},
				},
			}),
		}, nil
	}

	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("done", nil),
	}, nil
}

type multiRoundToolProvider struct {
	calls int
}

func (p *multiRoundToolProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	p.calls++
	switch p.calls {
	case 1:
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
				{
					ID:   "call-echo",
					Name: "echo",
					Arguments: map[string]any{
						"text": "first observation",
					},
				},
			}),
		}, nil
	case 2:
		transcript := toolTranscript(request.Messages)
		if !strings.Contains(transcript, "call-echo=first observation") {
			return llm.GenerateResponse{}, fmt.Errorf("missing first observation in %q", transcript)
		}
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
				{
					ID:   "call-calc",
					Name: "calculator",
					Arguments: map[string]any{
						"operation": "add",
						"a":         20,
						"b":         22,
					},
				},
			}),
		}, nil
	case 3:
		transcript := toolTranscript(request.Messages)
		if !strings.Contains(transcript, "call-calc=42") {
			return llm.GenerateResponse{}, fmt.Errorf("missing second observation in %q", transcript)
		}
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("multi-round complete", nil),
		}, nil
	default:
		return llm.GenerateResponse{}, fmt.Errorf("unexpected provider call %d", p.calls)
	}
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

func countToolResults(messages []conversation.Message) int {
	var count int
	for _, message := range messages {
		if _, ok := message.(conversation.ToolResultMessage); ok {
			count++
		}
	}
	return count
}

func newTestPromptBuilder() *prompt.Builder {
	return prompt.NewBuilder(prompt.SystemPrompt{
		SessionText: "system",
	})
}
