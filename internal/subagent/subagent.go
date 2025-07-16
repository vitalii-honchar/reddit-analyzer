package subagent

import (
	"context"
	"reddit-analyzer/internal/lib/tools"

	"github.com/sirupsen/logrus"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/agent"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

type SubAgentConfig struct {
	LLM      llm.LLMConfig
	Log      *logrus.Logger
	Tools    map[string]int
	Behavior string
}

func NewSubAgent[T any](cfg SubAgentConfig) (*agent.Agent[T], error) {
	tools, err := createTools[T](cfg.Tools)
	if err != nil {
		return nil, err
	}

	options := []agent.AgentOption[T]{
		agent.WithName[T]("search_agent"),
		agent.WithLLMConfig[T](cfg.LLM),
		agent.WithMiddleware[T](loggingMiddleware(cfg.Log)),
		agent.WithBehavior[T](cfg.Behavior),
	}
	options = append(options, tools...)

	return agent.NewAgent(options...)
}

func createTools[T any](tool map[string]int) ([]agent.AgentOption[T], error) {
	var options []agent.AgentOption[T]
	for name, limit := range tool {
		if name == tools.RedditSearchToolName {
			redditSearchTool, err := tools.NewRedditPostSearchTool()
			if err != nil {
				return nil, err
			}
			options = append(options, agent.WithTool[T](name, redditSearchTool))
		}
		options = append(options, agent.WithToolLimit[T](name, limit))
	}
	return options, nil
}

func loggingMiddleware(log *logrus.Logger) agent.AgentMiddleware {
	return func(ctx context.Context, state *agent.AgentState, msg llm.LLMMessage) (llm.LLMMessage, error) {

		log.WithFields(logrus.Fields{
			"message_type":   msg.Type,
			"agent_end":      msg.End,
			"tool_calls":     msg.ToolCalls,
			"content":        msg.Content[:64],
			"state_messages": len(state.Messages),
		}).Info("Received message from LLM")

		return msg, nil
	}
}
