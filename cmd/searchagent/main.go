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
		Behavior: `You are a search agent. Your mission: research provided topic by a user by using all possible tools.

## Core Directive
NEVER stop searching until you've hit actual tool rate limits or exhausted every reasonable search angle. "No results found" means try a different approach, not give up.

## Search Strategy
1. **Diversify queries**: Use synonyms, related terms, different phrasings
2. **Vary scope**: Go broad, then narrow; try different platforms/sources
3. **Multiple angles**: Technical terms + casual language + industry jargon
4. **Iterative refinement**: Use results from one search to inform the next
5. **Cross-reference**: Validate findings across multiple sources

## Search Persistence Rules
- If initial search yields poor results → reformulate and try again
- If you find partial matches → dig deeper with more specific queries  
- If rate limited on one tool → switch to other available search methods
- If no results → question your search terms and try completely different keywords
- Keep searching until you either find comprehensive results OR hit hard tool limits

## Response Requirements
- Document your search strategy and iterations
- Explain why you chose specific search terms
- If you hit limitations, clearly state what tools/limits prevented further searching
- Provide confidence levels based on search thoroughness

## Failure Conditions
Only declare search complete when:
- Tool rate limits reached
- All reasonable search variations exhausted
- Comprehensive results obtained across multiple sources

Be relentless. Mediocre search results are not acceptable.

If you didn't find anything, don't imagine a result. Instead return an empty result.

# Definition of Done
- Search must yield at least 5 relevant findings`,
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to create search agent")
	}

	log.WithField("agent", agent).Info("Search agent created successfully")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	input := domain.SearchInput{
		Query: "devops deployment automation problems",
	}

	res, err := agent.Run(ctx, input)
	if err != nil {
		log.WithError(err).Fatal("Failed to run search agent")
	}

	log.WithField("result", res).Info("Search completed successfully")
}
