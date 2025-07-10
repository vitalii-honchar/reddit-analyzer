package agent

import (
	"context"
	"fmt"
	"os"

	"reddit-analyzer/internal/redditanalyzer/domain"
	"reddit-analyzer/internal/redditanalyzer/tools"

	"github.com/sirupsen/logrus"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/agent"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

// RedditResearchAgent is the master ReAct agent that orchestrates the entire Reddit analysis workflow
type RedditResearchAgent struct {
	agent  agent.Agent[domain.AnalysisResult]
	logger *logrus.Logger
}

// NewRedditResearchAgent creates a new Reddit research ReAct agent with all required tools
func NewRedditResearchAgent(logger *logrus.Logger) (*RedditResearchAgent, error) {
	// Initialize OpenAI API key
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

	// Create the ReAct agent with proper configuration
	redditAgent, err := agent.NewAgent(
		agent.WithName[domain.AnalysisResult]("reddit_research_agent"),
		agent.WithLLMConfig[domain.AnalysisResult](llm.LLMConfig{
			Type:        llm.LLMTypeOpenAI,
			APIKey:      apiKey,
			Model:       "gpt-4o-mini",
			Temperature: 0.3,
		}),
		agent.WithBehavior[domain.AnalysisResult](createAgentBehavior()),
		
		// Tool registrations
		agent.WithTool[domain.AnalysisResult]("select_subreddits", subredditSelector.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("fetch_reddit_posts", postsFetcher.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("filter_posts", postsFilter.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("evaluate_post", postEvaluator.CreateLLMTool()),
		
		// Tool limit configurations
		agent.WithDefaultToolLimit[domain.AnalysisResult](5),  // Increase default from 3 to 5
		agent.WithToolLimit[domain.AnalysisResult]("select_subreddits", 1),   // Should only run once
		agent.WithToolLimit[domain.AnalysisResult]("fetch_reddit_posts", 1),  // Should only run once
		agent.WithToolLimit[domain.AnalysisResult]("filter_posts", 1),        // Should only run once
		agent.WithToolLimit[domain.AnalysisResult]("evaluate_post", 10),      // May need multiple calls for different posts
		
		)
	if err != nil {
		return nil, fmt.Errorf("failed to create ReAct agent: %w", err)
	}

	return &RedditResearchAgent{
		agent:  *redditAgent,
		logger: logger,
	}, nil
}

// createAgentBehavior defines the agent's behavior and reasoning instructions
func createAgentBehavior() string {
	return `You are a Reddit Research Agent specialized in discovering hidden market opportunities for indie hackers.

Your mission: Analyze Reddit to identify underdeveloped markets with low competition but high demand based on the user's project direction.

REASONING PROCESS (Think-Act-Observe):
1. THINK: Analyze the project direction and plan your approach
2. ACT: Use tools to gather and process data step by step
3. OBSERVE: Review tool results and decide on next actions
4. REPEAT: Continue the cycle until you have comprehensive analysis

WORKFLOW STEPS:
1. Use 'select_subreddits' ONCE to intelligently choose 3-5 niche subreddits
   - AVOID obvious entrepreneurship subreddits (r/entrepreneur, r/startups)
   - TARGET niche communities where problems exist but solutions are underdeveloped
   - FOCUS on domain-specific communities where professionals discuss pain points

2. Use 'fetch_reddit_posts' ONCE to get recent posts from selected subreddits (last 7 days)

3. Use 'filter_posts' ONCE to filter posts by engagement metrics (minimum upvotes/comments)

4. Use 'evaluate_post' for TOP 3-5 posts only to analyze for indie hacker opportunities
   - Rate each post 1-5 for opportunity potential
   - Focus on genuine pain points mentioned by domain professionals
   - Identify market gaps for affordable solutions

CRITICAL RULES:
- Do NOT generate business ideas - only analyze existing problems and trends
- Focus on problems that are underserved, not oversaturated markets
- Look for patterns in complaints, frustrations, and workflow inefficiencies
- Prioritize problems mentioned by professionals in their work contexts

When you receive a project direction, follow this EXACT workflow efficiently:
1. THINK about the project direction and plan your approach
2. ACT by using the available tools in this EXACT order:
   - select_subreddits: Call with {"project_direction": "the user's project direction"}
   - fetch_reddit_posts: Call with the selected subreddits
   - filter_posts: Call with the fetched posts
   - evaluate_post: Call with {"post": post_object, "project_direction": "the user's project direction"} for ONLY the top 3-5 posts
3. OBSERVE the results and compile your final analysis
4. STOP after completing step 3 - do not repeat the cycle

IMPORTANT: Always pass the original project direction to tools that require it (select_subreddits and evaluate_post).

Focus on genuine pain points mentioned by domain professionals, not oversaturated markets.`
}

// AnalyzeProject performs the complete Reddit analysis workflow using ReAct agent
func (r *RedditResearchAgent) AnalyzeProject(ctx context.Context, projectDirection string) (*domain.AnalysisResult, error) {
	r.logger.WithField("project_direction", projectDirection).Info("Starting ReAct agent analysis")

	// Run the ReAct agent with just the project direction
	r.logger.Info("🤖 Starting ReAct agent workflow...")
	agentResult, err := r.agent.Run(ctx, projectDirection)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ReAct agent workflow: %w", err)
	}

	r.logger.WithField("agent_result", agentResult.Data).Info("✅ ReAct agent analysis completed")
	return agentResult.Data, nil
}
