package react

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

func TestReActAgentRunsProvider(t *testing.T) {
	agent := NewReActAgent(
		answerTestProvider{},
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(nil)

	message := runAgentForAnswer(t, agent, session, "xin chao")

	if got, want := conversation.Text(message), "received: xin chao"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
}

func TestReActAgentSendsSystemPromptWithoutStoringItInConversation(t *testing.T) {
	provider := &captureMessagesProvider{}
	agent := NewReActAgent(
		provider,
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(nil)

	_ = runAgentForAnswer(t, agent, session, "xin chao")

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
	agent := NewReActAgent(
		usageProvider{},
		WithLogger(logger),
	)
	session := newTestSession(nil)

	_ = runAgentForAnswer(t, agent, session, "track stats")

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
	registry, err := tool.NewBaseRegistry(localtool.NewCalculator())
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	agent := NewReActAgent(
		&timelineProvider{},
		WithLogger(logger),
	)
	session := newTestSession(newTestToolSet(registry))

	_ = runAgentForAnswer(t, agent, session, "track timeline")

	output := logs.String()
	expectedParts := []string{
		`msg="Round 1: Think"`,
		`msg="Round 1: Act"`,
		"decision_summary=",
		"Need calculator to produce the requested value",
		`msg="Round 1: Tool calculator started"`,
		"arguments_preview=",
		`msg="Round 1: Observe calculator"`,
		"result_preview=42",
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
	registry, err := tool.NewBaseRegistry(errorTool{})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	provider := &toolErrorProvider{}
	agent := NewReActAgent(
		provider,
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(newTestToolSet(registry))

	message := runAgentForAnswer(t, agent, session, "test failing tool")

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
	registry, err := tool.NewBaseRegistry(blockingTool{})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	provider := &toolTimeoutProvider{}
	agent := NewReActAgent(
		provider,
		WithLogger(logging.NewNopLogger()),
		WithToolTimeout(time.Millisecond),
	)
	session := newTestSession(newTestToolSet(registry))

	message := runAgentForAnswer(t, agent, session, "test slow tool")

	if got, want := conversation.Text(message), "handled timeout"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if !strings.Contains(provider.observation, "tool error: execution timed out") {
		t.Fatalf("timeout observation missing: %q", provider.observation)
	}
}

func TestReActAgentExecutesToolCalls(t *testing.T) {
	registry, err := tool.NewBaseRegistry(localtool.NewCalculator())
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	provider := &scriptedToolProvider{}
	agent := NewReActAgent(
		provider,
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(newTestToolSet(registry))

	message := runAgentForAnswer(t, agent, session, "test tools")

	if got, want := conversation.Text(message), "all tools completed"; got != want {
		t.Fatalf("message text = %q, want %q", got, want)
	}
	if got, want := provider.calls, 2; got != want {
		t.Fatalf("provider calls = %d, want %d", got, want)
	}
}

func TestReActAgentRunsMultipleToolRounds(t *testing.T) {
	registry, err := tool.NewBaseRegistry(localtool.NewCalculator())
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	provider := &multiRoundToolProvider{}
	agent := NewReActAgent(
		provider,
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(newTestToolSet(registry))

	message := runAgentForAnswer(t, agent, session, "multi round")

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

func TestReActAgentEmitsToolEvents(t *testing.T) {
	registry, err := tool.NewBaseRegistry(localtool.NewCalculator())
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	runtimeAgent := NewReActAgent(
		&timelineProvider{},
		WithLogger(logging.NewNopLogger()),
	)
	session := newTestSession(newTestToolSet(registry))

	events := collectAgentEvents(t, runtimeAgent.Run(context.Background(), session, "track events"))
	assertEventType(t, events, agent.RunStartedEvent{})
	assertEventType(t, events, agent.RoundStartedEvent{})
	assertEventType(t, events, agent.ModelCallStartedEvent{})
	assertEventType(t, events, agent.ModelCallCompletedEvent{})
	assertEventType(t, events, agent.ToolCallStartedEvent{})
	assertEventType(t, events, agent.ToolCallCompletedEvent{})
	assertEventType(t, events, agent.FinalAnswerEvent{})

	completed := lastCompletedEvent(t, events)
	if completed.Reason != agent.TerminalCompleted {
		t.Fatalf("terminal reason = %q, want %q", completed.Reason, agent.TerminalCompleted)
	}
}

func TestReActAgentEmitsMaxRounds(t *testing.T) {
	registry, err := tool.NewBaseRegistry(localtool.NewCalculator())
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}

	runtimeAgent := NewReActAgent(
		&timelineProvider{},
		WithLogger(logging.NewNopLogger()),
		WithMaxRounds(1),
	)
	session := newTestSession(newTestToolSet(registry))

	events := collectAgentEvents(t, runtimeAgent.Run(context.Background(), session, "hit max rounds"))
	completed := lastCompletedEvent(t, events)
	if completed.Reason != agent.TerminalMaxRounds {
		t.Fatalf("terminal reason = %q, want %q", completed.Reason, agent.TerminalMaxRounds)
	}
}

func runAgentForAnswer(t *testing.T, runtimeAgent *ReActAgent, session *session.State, input string) conversation.Message {
	t.Helper()

	events := collectAgentEvents(t, runtimeAgent.Run(context.Background(), session, input))
	for _, event := range events {
		if failed, ok := event.(agent.RunFailedEvent); ok {
			t.Fatalf("Run() failed: %v", failed.Err)
		}
	}
	for i := len(events) - 1; i >= 0; i-- {
		if answer, ok := events[i].(agent.FinalAnswerEvent); ok {
			return answer.Message
		}
	}
	t.Fatalf("Run() emitted no final answer: %#v", events)
	return nil
}

func collectAgentEvents(t *testing.T, stream <-chan agent.Event) []agent.Event {
	t.Helper()

	events := []agent.Event{}
	for event := range stream {
		events = append(events, event)
	}
	return events
}

func assertEventType(t *testing.T, events []agent.Event, target agent.Event) {
	t.Helper()

	targetType := reflect.TypeOf(target)
	for _, event := range events {
		if reflect.TypeOf(event) == targetType {
			return
		}
	}
	t.Fatalf("missing event type %T in %#v", target, events)
}

func lastCompletedEvent(t *testing.T, events []agent.Event) agent.RunCompletedEvent {
	t.Helper()

	for i := len(events) - 1; i >= 0; i-- {
		if completed, ok := events[i].(agent.RunCompletedEvent); ok {
			return completed
		}
	}
	t.Fatalf("missing agent.RunCompletedEvent in %#v", events)
	return agent.RunCompletedEvent{}
}

type scriptedToolProvider struct {
	calls int
}

type usageProvider struct{}

type answerTestProvider struct{}

func (p answerTestProvider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
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
		Usage: usage.TokenUsage{
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
			Message: conversation.NewAssistantMessage("Need calculator to produce the requested value.", []conversation.ToolCall{
				{
					ID:   "call-calc",
					Name: "calculator",
					Arguments: map[string]any{
						"operation": "multiply",
						"a":         6,
						"b":         7,
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
					ID:   "call-first",
					Name: "calculator",
					Arguments: map[string]any{
						"operation": "add",
						"a":         1,
						"b":         1,
					},
				},
			}),
		}, nil
	case 2:
		transcript := toolTranscript(request.Messages)
		if !strings.Contains(transcript, "call-first=2") {
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
		if got, want := len(request.Tools), 1; got != want {
			return llm.GenerateResponse{}, fmt.Errorf("tool count = %d, want %d", got, want)
		}
		return llm.GenerateResponse{
			Message: conversation.NewAssistantMessage("", []conversation.ToolCall{
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
					ID:   "call-add",
					Name: "calculator",
					Arguments: map[string]any{
						"operation": "add",
						"a":         20,
						"b":         22,
					},
				},
			}),
		}, nil
	}

	transcript := toolTranscript(request.Messages)
	if !strings.Contains(transcript, "call-calc=42") {
		return llm.GenerateResponse{}, fmt.Errorf("missing calculator result in %q", transcript)
	}
	if !strings.Contains(transcript, "call-add=42") {
		return llm.GenerateResponse{}, fmt.Errorf("missing second calculator result in %q", transcript)
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
		InputSchema: tool.EmptyInputSchema(),
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
		InputSchema: tool.EmptyInputSchema(),
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

func newTestToolSet(registry *tool.BaseRegistry) *tool.ToolSet {
	return tool.NewToolSet(registry, tool.NewRuntimeRegistry(20))
}

func newTestSession(tools *tool.ToolSet) *session.State {
	return session.NewState(
		session.WithPromptSnapshot(prompt.SessionPrompt{
			DynamicParts: []prompt.DynamicPart{"system"},
		}),
		session.WithTools(tools),
	)
}
