package agent

type ToolName string

type Tool struct {
	Name             ToolName                                       `json:"name"`
	ParametersSchema map[string]any                                 `json:"parameters_schema"`
	Description      string                                         `json:"description"`
	Call             func(args map[string]any) (*ToolResult, error) `json:"-"`
}

type ToolResult interface {
	Json() string
}

type ToolsLimits struct {
	Limits map[ToolName]int `json:"limits"`
}

type ToolCall struct {
	ToolName ToolName       `json:"tool_name"`
	Args     map[string]any `json:"args"`
}

func NewToolCall(toolName ToolName, args map[string]any) *ToolCall {
	return &ToolCall{
		ToolName: toolName,
		Args:     args,
	}
}

func isLimitReached(usage map[ToolName]int, limits map[ToolName]int) bool {
	limitReached := make(map[ToolName]bool)
	for toolName, limit := range limits {
		if usage, exists := usage[toolName]; exists && usage >= limit {
			limitReached[toolName] = true
		}
	}
	return len(limitReached) == len(limits)
}
