package llm

type LLMTool struct {
	Name             string         `json:"name"`
	ParametersSchema map[string]any `json:"parameters_schema"`
	Description      string         `json:"description"`
}

type LLMToolResult interface {
	Json() string
}

type LLMToolCall struct {
	ToolName string         `json:"tool_name"`
	Args     map[string]any `json:"args"`
}

func NewLLMToolCall(toolName string, args map[string]any) LLMToolCall {
	return LLMToolCall{
		ToolName: toolName,
		Args:     args,
	}
}
