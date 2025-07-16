package main

import (
	"context"
	"os"
	"reddit-analyzer/internal/lib/domain"
	"reddit-analyzer/internal/lib/tools"
	"reddit-analyzer/internal/subagent"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	agent, err := subagent.NewSubAgent[domain.SearchResult](subagent.SubAgentConfig{
		LLM: llm.LLMConfig{
			Type:        llm.LLMTypeOpenAI,
			APIKey:      apiKey,
			Model:       "gpt-4.1",
			Temperature: 0.0,
		},
		Log: log,
		Tools: map[string]int{
			tools.RedditSearchToolName: 50,
		},
		Behavior: "You are a Reddit search agent that helps find relevant posts based on user queries.",
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to create search agent")
	}

	log.WithField("agent", agent).Info("Search agent created successfully")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	input := domain.SearchInput{
		Query: "Find subreddits where backend engineers complain about deployment automation",
	}

	res, err := agent.Run(ctx, input)
	if err != nil {
		log.WithError(err).Fatal("Failed to run search agent")
	}

	log.WithField("result", res).Info("Search completed successfully")
}
