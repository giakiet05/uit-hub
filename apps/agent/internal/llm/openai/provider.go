package openai

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	openaisdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type Config struct {
	APIKey  string
	Model   string
	BaseURL string
}

type Provider struct {
	model  string
	client openaisdk.Client
	logger *slog.Logger
}

func NewProvider(cfg Config, logger *slog.Logger) *Provider {
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	options := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
	}
	if cfg.BaseURL != "" {
		options = append(options, option.WithBaseURL(strings.TrimRight(cfg.BaseURL, "/")))
	}

	return &Provider{
		model:  cfg.Model,
		client: openaisdk.NewClient(options...),
		logger: logger,
	}
}

func (p *Provider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	modelName := request.Model
	if modelName == "" {
		modelName = p.model
	}

	p.logger.DebugContext(ctx, "Calling OpenAI Responses API", "model", modelName, "tool_count", len(request.Tools))
	response, err := p.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: shared.ResponsesModel(modelName),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toResponsesInput(request.Messages),
		},
		Tools: toResponsesTools(request.Tools),
	})
	if err != nil {
		p.logger.DebugContext(ctx, "OpenAI Responses API call failed", "model", modelName, "error", err)
		return llm.GenerateResponse{}, err
	}
	p.logger.DebugContext(ctx, "OpenAI Responses API call completed", "model", modelName, "input_tokens", response.Usage.InputTokens, "output_tokens", response.Usage.OutputTokens)

	message, err := toLLMMessage(response.OutputText(), response.Output)
	if err != nil {
		p.logger.DebugContext(ctx, "OpenAI Responses API output mapping failed", "model", modelName, "error", err)
		return llm.GenerateResponse{}, err
	}

	return llm.GenerateResponse{
		Message: message,
		Usage: llm.Usage{
			InputTokens:  int(response.Usage.InputTokens),
			OutputTokens: int(response.Usage.OutputTokens),
		},
	}, nil
}

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

func toResponsesFunctionCall(call conversation.ToolCall) responses.ResponseInputItemUnionParam {
	arguments, err := json.Marshal(call.Arguments)
	if err != nil {
		arguments = []byte("{}")
	}
	return responses.ResponseInputItemParamOfFunctionCall(string(arguments), call.ID, call.Name)
}

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
				Parameters:  tool.InputSchema,
				Strict:      openaisdk.Bool(false),
			},
		})
	}
	return params
}

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
