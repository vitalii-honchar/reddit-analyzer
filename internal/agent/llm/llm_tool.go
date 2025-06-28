package llm

type LLMTool[T LLMToolResult] struct {
	Name             string                                          `json:"name"`
	ParametersSchema map[string]any                                  `json:"parameters_schema"`
	Description      string                                          `json:"description"`
	Call             func(id string, args map[string]any) (T, error) `json:"-"`
}

type LLMToolResult struct {
	ID string `json:"id"`
}

type LLMToolCall struct {
	ID       string         `json:"id"`
	ToolName string         `json:"tool_name"`
	Args     map[string]any `json:"args"`
}

func NewLLMToolCall(id string, toolName string, args map[string]any) LLMToolCall {
	return LLMToolCall{
		ID:       id,
		ToolName: toolName,
		Args:     args,
	}
}
