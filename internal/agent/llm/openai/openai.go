package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"reddit-analyzer/internal/agent/agent"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAILLM struct {
	client      openai.Client
	apiKey      string
	temperature float64
	model       openai.ChatModel
	tools       []agent.Tool
}

type OpenAILLMOption func(o *OpenAILLM)

func WithOpenAILLMTemperature(temperature float64) OpenAILLMOption {
	return func(o *OpenAILLM) {
		o.temperature = temperature
	}
}
func WithOpenAILLMModel(model openai.ChatModel) OpenAILLMOption {
	return func(o *OpenAILLM) {
		o.model = model
	}
}

func WithAPIKey(apiKey string) OpenAILLMOption {
	return func(o *OpenAILLM) {
		o.apiKey = apiKey
		o.client = openai.NewClient(option.WithAPIKey(apiKey))
	}
}

func WithTools(tools []agent.Tool) OpenAILLMOption {
	return func(o *OpenAILLM) {
		o.tools = tools
	}

}

func NewOpenAILLM(options ...OpenAILLMOption) *OpenAILLM {
	llm := &OpenAILLM{}
	for _, opt := range options {
		opt(llm)
	}
	return llm
}

func (o *OpenAILLM) Call(ctx context.Context, msgs []agent.LLMMessage) (agent.LLMMessage, error) {
	// Convert agent messages to OpenAI format
	params := o.createParameters(msgs)

	// Make API call
	completion, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return agent.LLMMessage{}, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return agent.LLMMessage{}, fmt.Errorf("no response from OpenAI")
	}

	choice := completion.Choices[0]
	message := choice.Message

	// Convert response back to agent format
	result := agent.LLMMessage{
		Type:    agent.LLMMessageTypeAssistant,
		Content: message.Content,
	}

	// Handle tool calls in response
	if len(message.ToolCalls) > 0 {
		toolCall := message.ToolCalls[0] // Take first tool call
		args := make(map[string]any)
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			return agent.LLMMessage{}, fmt.Errorf("failed to parse tool arguments: %w", err)
		}

		result.ToolCall = agent.NewToolCall(agent.ToolName(toolCall.Function.Name), args)
	}

	if choice.FinishReason == "stop" || choice.FinishReason == "length" {
		result.End = true
	}

	return result, nil
}

func (o *OpenAILLM) createParameters(messages []agent.LLMMessage) openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Messages:    o.createMessages(messages),
		Model:       o.model,
		Temperature: openai.Float(o.temperature),
		Tools:       o.createToolParams(),
	}
}

func (o *OpenAILLM) createToolParams() []openai.ChatCompletionToolParam {
	toolParams := make([]openai.ChatCompletionToolParam, 0, len(o.tools))

	for _, tool := range o.tools {
		toolParams = append(toolParams, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        string(tool.Name),
				Description: openai.String(tool.Description),
				Parameters:  tool.ParametersSchema,
			},
		})
	}

	return toolParams
}

func (o *OpenAILLM) createMessages(msgs []agent.LLMMessage) []openai.ChatCompletionMessageParamUnion {
	openAIMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(msgs))

	for _, msg := range msgs {
		switch msg.Type {
		case agent.LLMMessageTypeSystem:
			openAIMessages = append(openAIMessages, openai.SystemMessage(msg.Content))
		case agent.LLMMessageTypeUser:
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		case agent.LLMMessageTypeAssistant:
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		}
	}

	return openAIMessages
}

func (o *OpenAILLM) marshalArgs(args map[string]any) string {
	if args == nil {
		return "{}"
	}

	bytes, err := json.Marshal(args)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
