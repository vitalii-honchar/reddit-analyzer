package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
	"reddit-analyzer/internal/redditanalyzer/domain"
	"reddit-analyzer/internal/redditanalyzer/tools"
)

// RedditResearchAgent is the master agent that orchestrates the entire Reddit analysis workflow
type RedditResearchAgent struct {
	llm    llm.LLM
	logger *logrus.Logger
}

// NewRedditResearchAgent creates a new Reddit research agent with all required tools
func NewRedditResearchAgent(logger *logrus.Logger) (*RedditResearchAgent, error) {
	// Initialize OpenAI LLM
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	// Create tool instances
	subredditSelector, err := tools.NewSubredditSelectorTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create subreddit selector: %w", err)
	}

	postsFetcher := tools.NewTopSubredditPostsTool()
	postsFilter := tools.NewFilterPostsTool()

	postEvaluator, err := tools.NewEvaluatePostTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create post evaluator: %w", err)
	}

	// Create tools map for LLM
	toolsMap := map[string]llm.LLMTool{
		"select_subreddits":   subredditSelector.CreateLLMTool(),
		"fetch_reddit_posts":  postsFetcher.CreateLLMTool(),
		"filter_posts":        postsFilter.CreateLLMTool(),
		"evaluate_post":       postEvaluator.CreateLLMTool(),
	}

	// Create LLM with tools
	openaiLLM, err := llm.CreateLLM(llm.LLMConfig{
		Type:        llm.LLMTypeOpenAI,
		APIKey:      apiKey,
		Model:       "gpt-4o",
		Temperature: 0.7,
	}, toolsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI LLM: %w", err)
	}

	return &RedditResearchAgent{
		llm:    openaiLLM,
		logger: logger,
	}, nil
}

// AnalyzeProject performs the complete Reddit analysis workflow using LLM agent orchestration
func (r *RedditResearchAgent) AnalyzeProject(ctx context.Context, projectDirection string) (*domain.AnalysisResult, error) {
	r.logger.WithField("project_direction", projectDirection).Info("Starting Reddit analysis")

	// Create the analysis prompt for the LLM agent
	analysisPrompt := fmt.Sprintf(`
You are a Reddit Research Agent for indie hackers to discover hidden market opportunities. 
Analyze Reddit to identify underdeveloped markets with low competition but high demand.

Project Direction: "%s"

Follow this workflow:
1. Use select_subreddits tool to intelligently select 3-5 niche subreddits based on the project direction
   - Avoid obvious entrepreneurship subreddits (r/entrepreneur, r/startups)
   - Target niche communities where problems exist but solutions are underdeveloped
   
2. Use fetch_reddit_posts tool to get recent posts from selected subreddits (last 7 days)

3. Use filter_posts tool to filter posts by engagement metrics (minimum upvotes/comments)

4. Use evaluate_post tool to analyze each filtered post for hidden indie hacker opportunities
   - Rate each post 1-5 for opportunity potential
   - Focus on genuine pain points mentioned by domain professionals
   - Identify market gaps for affordable solutions

Please execute this workflow step by step and provide a comprehensive analysis of the opportunities found.
`, projectDirection)

	// Create LLM messages
	messages := []llm.LLMMessage{
		{
			Type:    llm.LLMMessageTypeSystem,
			Content: "You are a Reddit Research Agent specializing in finding hidden market opportunities for indie hackers.",
		},
		{
			Type:    llm.LLMMessageTypeUser,
			Content: analysisPrompt,
		},
	}

	// Call the LLM with tools
	r.logger.Info("🤖 Starting LLM-powered analysis workflow...")
	response, err := r.llm.Call(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to execute LLM workflow: %w", err)
	}

	r.logger.WithField("response", response.Content).Info("✅ Analysis workflow completed")

	// For now, return a basic result - in a real implementation, you'd parse the LLM response
	// and extract the structured data from tool results
	result := &domain.AnalysisResult{
		ProjectDirection:   projectDirection,
		SelectedSubreddits: []string{"placeholder"}, // Would be extracted from tool results
		PostsAnalyzed:      0,                       // Would be extracted from tool results
		PostsFiltered:      0,                       // Would be extracted from tool results
		Opportunities:      []domain.RankedOpportunity{}, // Would be extracted from tool results
		AnalysisTime:       time.Now(),
	}

	return result, nil
}