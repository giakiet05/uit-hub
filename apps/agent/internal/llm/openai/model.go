package openai

import (
	"context"
	"encoding/json"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
	openaisdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

// Model implements llm.StreamModel for OpenAI.
type Model struct {
	llm.BaseModel
	client openaisdk.Client
}

// Generate sends conversation history and tool definitions to OpenAI and maps
// the response back into internal conversation types.
func (m *Model) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	modelName := m.Name()

	m.Logger().DebugContext(ctx, "Calling OpenAI Responses API", "model", modelName, "tool_count", len(request.Tools))
	response, err := m.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: shared.ResponsesModel(modelName),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toResponsesInput(request.Messages),
		},
		Tools: toResponsesTools(request.Tools),
	})
	if err != nil {
		m.Logger().DebugContext(ctx, "OpenAI Responses API call failed", "model", modelName, "error", err)
		return llm.GenerateResponse{}, err
	}
	m.Logger().DebugContext(ctx, "OpenAI Responses API call completed", "model", modelName, "input_tokens", response.Usage.InputTokens, "output_tokens", response.Usage.OutputTokens)

	message, err := toLLMMessage(response.OutputText(), response.Output)
	if err != nil {
		m.Logger().DebugContext(ctx, "OpenAI Responses API output mapping failed", "model", modelName, "error", err)
		return llm.GenerateResponse{}, err
	}

	return llm.GenerateResponse{
		Message: message,
		Usage: usage.TokenUsage{
			InputTokens:  int(response.Usage.InputTokens),
			OutputTokens: int(response.Usage.OutputTokens),
		},
	}, nil
}

// Stream sends a streaming Responses API request and emits visible assistant
// text deltas before the final mapped response.
func (m *Model) Stream(ctx context.Context, request llm.GenerateRequest) <-chan llm.StreamEvent {
	events := make(chan llm.StreamEvent)
	go func() {
		defer close(events)

		modelName := m.Name()

		m.Logger().DebugContext(ctx, "Calling OpenAI Responses streaming API", "model", modelName, "tool_count", len(request.Tools))
		stream := m.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
			Model: shared.ResponsesModel(modelName),
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: toResponsesInput(request.Messages),
			},
			Tools: toResponsesTools(request.Tools),
		})
		defer stream.Close()

		for stream.Next() {
			event := stream.Current()
			switch typed := event.AsAny().(type) {
			case responses.ResponseTextDeltaEvent:
				if typed.Delta != "" && !emitStreamEvent(ctx, events, llm.TextDeltaEvent{Delta: typed.Delta}) {
					return
				}
			case responses.ResponseCompletedEvent:
				message, err := toLLMMessage(typed.Response.OutputText(), typed.Response.Output)
				if err != nil {
					emitStreamEvent(ctx, events, llm.StreamFailedEvent{Err: err})
					return
				}

				m.Logger().DebugContext(ctx, "OpenAI Responses streaming API call completed", "model", modelName, "input_tokens", typed.Response.Usage.InputTokens, "output_tokens", typed.Response.Usage.OutputTokens)
				emitStreamEvent(ctx, events, llm.StreamCompletedEvent{
					Response: llm.GenerateResponse{
						Message: message,
						Usage: usage.TokenUsage{
							InputTokens:  int(typed.Response.Usage.InputTokens),
							OutputTokens: int(typed.Response.Usage.OutputTokens),
						},
					},
				})
				return
			}
		}

		if err := stream.Err(); err != nil {
			m.Logger().DebugContext(ctx, "OpenAI Responses streaming API call failed", "model", modelName, "error", err)
			emitStreamEvent(ctx, events, llm.StreamFailedEvent{Err: err})
		}
	}()
	return events
}

func emitStreamEvent(ctx context.Context, events chan<- llm.StreamEvent, event llm.StreamEvent) bool {
	select {
	case <-ctx.Done():
		return false
	case events <- event:
		return true
	}
}

// toResponsesInput converts internal messages to OpenAI Responses input items.
func toResponsesInput(messages []conversation.Message) []responses.ResponseInputItemUnionParam {
	input := make([]responses.ResponseInputItemUnionParam, 0, len(messages))
	for _, message := range messages {
		switch typed := message.(type) {
		case conversation.SystemMessage, conversation.UserMessage:
			if text := conversation.Text(message); text != "" {
				input = append(input, responses.ResponseInputItemParamOfMessage(text, toResponsesRole(conversation.RoleOf(message))))
			}
		case conversation.AssistantMessage:
			if text := conversation.Text(typed); text != "" {
				input = append(input, responses.ResponseInputItemParamOfMessage(text, responses.EasyInputMessageRoleAssistant))
			}
			for _, call := range typed.ToolCalls {
				input = append(input, toResponsesFunctionCall(call))
			}
		case conversation.ToolResultMessage:
			input = append(input, responses.ResponseInputItemParamOfFunctionCallOutput(typed.ToolCallID, conversation.Text(typed)))
		}
	}
	return input
}

// toResponsesFunctionCall converts an internal tool call back to an OpenAI
// function-call input item for the next round.
func toResponsesFunctionCall(call conversation.ToolCall) responses.ResponseInputItemUnionParam {
	arguments, err := json.Marshal(call.Arguments)
	if err != nil {
		arguments = []byte("{}")
	}
	return responses.ResponseInputItemParamOfFunctionCall(string(arguments), call.ID, call.Name)
}

// toResponsesTools converts internal tool definitions to OpenAI function tools.
func toResponsesTools(tools []tool.Definition) []responses.ToolUnionParam {
	if len(tools) == 0 {
		return nil
	}

	params := make([]responses.ToolUnionParam, 0, len(tools))
	for _, tool := range tools {
		params = append(params, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.Name,
				Description: openaisdk.String(tool.Description),
				Parameters:  tool.InputSchema.Clone(),
				Strict:      openaisdk.Bool(false),
			},
		})
	}
	return params
}

// toLLMMessage maps OpenAI output text and function calls to an internal
// assistant message.
func toLLMMessage(outputText string, output []responses.ResponseOutputItemUnion) (conversation.AssistantMessage, error) {
	toolCalls := []conversation.ToolCall{}
	for _, item := range output {
		functionCall, ok := item.AsAny().(responses.ResponseFunctionToolCall)
		if !ok {
			continue
		}

		arguments := map[string]any{}
		if functionCall.Arguments != "" {
			if err := json.Unmarshal([]byte(functionCall.Arguments), &arguments); err != nil {
				return conversation.AssistantMessage{}, err
			}
		}

		toolCalls = append(toolCalls, conversation.ToolCall{
			ID:        functionCall.CallID,
			Name:      functionCall.Name,
			Arguments: arguments,
		})
	}
	return conversation.NewAssistantMessage(outputText, toolCalls), nil
}

// toResponsesRole maps internal message roles to OpenAI input roles.
func toResponsesRole(role conversation.Role) responses.EasyInputMessageRole {
	switch role {
	case conversation.RoleSystem:
		return responses.EasyInputMessageRoleSystem
	case conversation.RoleAssistant:
		return responses.EasyInputMessageRoleAssistant
	case conversation.RoleUser:
		return responses.EasyInputMessageRoleUser
	default:
		return responses.EasyInputMessageRoleUser
	}
}
