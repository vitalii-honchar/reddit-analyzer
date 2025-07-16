package searchagent

import (
	"context"
	"reddit-analyzer/internal/lib/tools"
	"reddit-analyzer/internal/searchagent/domain"

	"github.com/sirupsen/logrus"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/agent"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

const (
	redditSearchToolLimit = 10
)

type SearchAgentConfig struct {
	LLM      llm.LLMConfig
	Log      *logrus.Logger
	Behavior string
}

func NewSearchAgent(cfg SearchAgentConfig) (*agent.Agent[domain.SearchResult], error) {

	redditSearchTool, err := tools.NewRedditPostSearchTool()
	if err != nil {
		return nil, err
	}

	return agent.NewAgent(
		agent.WithName[domain.SearchResult]("search_agent"),
		agent.WithLLMConfig[domain.SearchResult](cfg.LLM),
		agent.WithMiddleware[domain.SearchResult](loggingMiddleware(cfg.Log)),
		agent.WithTool[domain.SearchResult](redditSearchTool.Name, redditSearchTool),
		agent.WithToolLimit[domain.SearchResult](redditSearchTool.Name, redditSearchToolLimit),
		agent.WithBehavior[domain.SearchResult](cfg.Behavior),
	)
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
