package agent

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ErrLimitReached        = errors.New("tool limit reached")
	ErrToolError           = errors.New("tool error occurred")
	ErrLLMCall             = errors.New("LLM call error occurred")
	ErrFinish              = errors.New("LLM finished execution")
	ErrToolNotFound        = errors.New("tool not found")
	ErrInvalidResultSchema = errors.New("invalid result schema")
	ErrCannotCreateSchema  = errors.New("cannot create schema from output type")
)

var systemPromptTemplate = NewPrompt(`
You are an agent that should act as specified in escaped content <BEHAVIOR></BEHAVIOR>.
At the end of execution when you will be read to finish, you should return a JSON object that matches the output schema.

TOOLS AVAILABLE TO USE:
{{.tools}}

TOOLS USAGE LIMITS:
{{.tools_usage}}

TOOLS CALLING LIMITS:
{{.calling_limits}}

OUTPUT SCHEMA:
{{.output_schema}}

<BEHAVIOR>
{{.behavior}}
</BEHAVIOR>
`)

type Agent[T any] struct {
	llm              LLM
	tools            map[ToolName]Tool
	limits           map[ToolName]int
	outputSchema     *T
	systemPrompt     Prompt
	behavior         string
	outputJSONSchema string
}

type AgentState struct {
	Messages []LLMMessage
}

func (a *AgentState) AddMessage(msg LLMMessage) {
	a.Messages = append(a.Messages, msg)
}

type LLM interface {
	Call(msgs []LLMMessage) (LLMMessage, error)
}

func (a *Agent[T]) Run() (*AgentResult[T], error) {

	state, err := a.createInitState()
	if err != nil {
		return nil, err
	}
	usage := make(map[ToolName]int)

	for {
		if isLimitReached(usage, a.limits) {
			res, err := a.createResult(state)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrLimitReached, err)
			}
			return res, ErrLimitReached
		}

		llmMessage, err := a.llm.Call(state.Messages)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrLLMCall, err)
		}

		if llmMessage.ToolCall != nil {
			if err := a.callTool(&llmMessage, usage); err != nil {
				return nil, fmt.Errorf("%w: %s", ErrToolError, err)
			}
		}

		state.AddMessage(llmMessage)

		if llmMessage.End {
			return a.createResult(state)
		}

		newSystemPrompt, err := a.createSystemPrompt(usage)
		if err != nil {
			return nil, fmt.Errorf("failed to update system prompt: %w", err)
		}
		state.Messages[0].Content = newSystemPrompt
	}
}

func (a *Agent[T]) createInitState() (*AgentState, error) {
	systemPrompt, err := a.createSystemPrompt(make(map[ToolName]int))
	if err != nil {
		return nil, fmt.Errorf("failed to create system prompt: %w", err)
	}
	return &AgentState{
		Messages: []LLMMessage{
			NewLLMMessage(LLMMessageTypeSystem, systemPrompt),
		},
	}, nil
}

func (a *Agent[T]) createSystemPrompt(usage map[ToolName]int) (string, error) {
	schema := jsonschema.Reflect(a.outputSchema)
	outputSchema, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}

	tools, err := json.Marshal(a.tools)
	if err != nil {
		return "", err
	}

	toolsUsage, err := json.Marshal(usage)
	if err != nil {
		return "", err
	}

	callingLimits, err := json.Marshal(a.limits)
	if err != nil {
		return "", err
	}

	return a.systemPrompt.Render(map[string]any{
		"tools":          string(tools),
		"tools_usage":    string(toolsUsage),
		"calling_limits": string(callingLimits),
		"output_schema":  string(outputSchema),
		"behavior":       a.behavior,
	})
}

func (a *Agent[T]) callTool(llmMessage *LLMMessage, usage map[ToolName]int) error {
	tool, ok := a.tools[llmMessage.ToolCall.ToolName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrToolNotFound, llmMessage.ToolCall.ToolName)
	}
	toolRes, err := tool.Call(llmMessage.ToolCall.Args)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrToolError, err)
	}
	usage[llmMessage.ToolCall.ToolName]++
	llmMessage.ToolResult = toolRes

	return nil
}

func (a *Agent[T]) createResult(state *AgentState) (*AgentResult[T], error) {
	dataLoader := gojsonschema.NewStringLoader(state.Messages[len(state.Messages)-1].Content)
	schemaLoader := gojsonschema.NewStringLoader(a.outputJSONSchema)
	validationRes, err := gojsonschema.Validate(schemaLoader, dataLoader)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResultSchema, err)
	}
	if !validationRes.Valid() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResultSchema, validationRes.Errors())
	}

	var data T

	if err := json.Unmarshal([]byte(state.Messages[len(state.Messages)-1].Content), &data); err != nil {
		return nil, err
	}

	return &AgentResult[T]{
		Data:     &data,
		Messages: state.Messages,
	}, nil
}
