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
	llm          LLM
	tools        map[ToolName]Tool
	limits       map[ToolName]int
	outputSchema *T
	systemPrompt Prompt
	behavior     string
	schemaLoader gojsonschema.JSONLoader
}

type AgentOption[T any] func(*Agent[T])

func NewAgent[T any](options ...AgentOption[T]) *Agent[T] {
	agent := &Agent[T]{
		tools:  make(map[ToolName]Tool),
		limits: make(map[ToolName]int),
	}
	for _, opt := range options {
		opt(agent)
	}
	return agent
}

// WithLLM sets the LLM for the agent
func WithLLM[T any](llm LLM) AgentOption[T] {
	return func(a *Agent[T]) {
		a.llm = llm
	}
}

// WithBehavior sets the behavior description for the agent
func WithBehavior[T any](behavior string) AgentOption[T] {
	return func(a *Agent[T]) {
		a.behavior = behavior
	}
}

// WithOutputSchema sets the output schema for type-safe result parsing
func WithOutputSchema[T any](schema *T) AgentOption[T] {
	return func(a *Agent[T]) {
		a.outputSchema = schema
		// Pre-compile schema for validation
		reflectedSchema := jsonschema.Reflect(schema)
		schemaBytes, err := json.Marshal(reflectedSchema)
		if err == nil {
			a.schemaLoader = gojsonschema.NewStringLoader(string(schemaBytes))
		}
	}
}

// WithSystemPrompt sets a custom system prompt template
func WithSystemPrompt[T any](prompt Prompt) AgentOption[T] {
	return func(a *Agent[T]) {
		a.systemPrompt = prompt
	}
}

// WithTool adds a tool to the agent
func WithTool[T any](name ToolName, tool Tool) AgentOption[T] {
	return func(a *Agent[T]) {
		a.tools[name] = tool
	}
}

// WithToolLimit sets the usage limit for a specific tool
func WithToolLimit[T any](name ToolName, limit int) AgentOption[T] {
	return func(a *Agent[T]) {
		a.limits[name] = limit
	}
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
	validationRes, err := gojsonschema.Validate(a.schemaLoader, dataLoader)

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
