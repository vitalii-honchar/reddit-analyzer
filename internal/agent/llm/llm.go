package llm

import (
	"context"
	"fmt"
)

var ErrUnsupportedLLMType = fmt.Errorf("unsupported LLM type")

type LLM interface {
	Call(ctx context.Context, msgs []LLMMessage) (LLMMessage, error)
}

func CreateLLM(cfg LLMConfig, tools []LLMTool) (LLM, error) {
	switch cfg.Type {
	case LLMTypeOpenAI:
		return newOpenAILLM(
			withOpenAIAPIKey(cfg.APIKey),
			withOpenAILLMModel(cfg.Model),
			withOpenAILLMTemperature(cfg.Temperature),
			withOpenAITools(tools),
		), nil
	default:
		return nil, ErrUnsupportedLLMType
	}
}
