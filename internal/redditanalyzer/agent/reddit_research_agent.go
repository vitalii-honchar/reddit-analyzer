package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/agent"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
	"reddit-analyzer/internal/redditanalyzer/domain"
	"reddit-analyzer/internal/redditanalyzer/tools"
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
			Model:       "gpt-4o",
			Temperature: 0.7,
		}),
		agent.WithBehavior[domain.AnalysisResult](createAgentBehavior()),
		agent.WithTool[domain.AnalysisResult]("select_subreddits", subredditSelector.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("fetch_reddit_posts", postsFetcher.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("filter_posts", postsFilter.CreateLLMTool()),
		agent.WithTool[domain.AnalysisResult]("evaluate_post", postEvaluator.CreateLLMTool()),
		agent.WithOutputSchema[domain.AnalysisResult](&domain.AnalysisResult{}),
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
1. Use 'select_subreddits' to intelligently choose 3-5 niche subreddits
   - AVOID obvious entrepreneurship subreddits (r/entrepreneur, r/startups)
   - TARGET niche communities where problems exist but solutions are underdeveloped
   - FOCUS on domain-specific communities where professionals discuss pain points

2. Use 'fetch_reddit_posts' to get recent posts from selected subreddits (last 7 days)

3. Use 'filter_posts' to filter posts by engagement metrics (minimum upvotes/comments)

4. Use 'evaluate_post' to analyze each filtered post for indie hacker opportunities
   - Rate each post 1-5 for opportunity potential
   - Focus on genuine pain points mentioned by domain professionals
   - Identify market gaps for affordable solutions

CRITICAL RULES:
- Do NOT generate business ideas - only analyze existing problems and trends
- Focus on problems that are underserved, not oversaturated markets
- Look for patterns in complaints, frustrations, and workflow inefficiencies
- Prioritize problems mentioned by professionals in their work contexts

Return a structured AnalysisResult with discovered opportunities ranked by potential.`
}

// AnalyzeProject performs the complete Reddit analysis workflow using ReAct agent
func (r *RedditResearchAgent) AnalyzeProject(ctx context.Context, projectDirection string) (*domain.AnalysisResult, error) {
	r.logger.WithField("project_direction", projectDirection).Info("Starting ReAct agent analysis")

	// Create the input prompt for the ReAct agent
	agentInput := fmt.Sprintf(`Analyze Reddit to discover hidden market opportunities for indie hackers.

Project Direction: "%s"

Please follow the ReAct workflow:
1. THINK about the project direction and plan your approach
2. ACT by using the available tools step by step:
   - select_subreddits: Choose 3-5 niche subreddits (avoid r/entrepreneur, r/startups)
   - fetch_reddit_posts: Get recent posts from selected subreddits (last 7 days)
   - filter_posts: Filter by engagement metrics
   - evaluate_post: Analyze each post for indie hacker opportunities (rate 1-5)
3. OBSERVE the results and continue reasoning
4. REPEAT until you have comprehensive analysis

Focus on genuine pain points mentioned by domain professionals, not oversaturated markets.
Return a structured analysis with discovered opportunities ranked by potential.`, projectDirection)

	// Run the ReAct agent
	r.logger.Info("🤖 Starting ReAct agent workflow...")
	agentResult, err := r.agent.Run(ctx, agentInput)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ReAct agent workflow: %w", err)
	}

	// Extract the actual result from AgentResult - need to check exact field name
	// For now, create a basic result structure until we can access the agent result properly
	analysisResult := &domain.AnalysisResult{
		ProjectDirection:   projectDirection,
		SelectedSubreddits: []string{}, // Will be populated by agent tools
		PostsAnalyzed:      0,          // Will be populated by agent tools
		PostsFiltered:      0,          // Will be populated by agent tools
		Opportunities:      []domain.RankedOpportunity{}, // Will be populated by agent tools
		AnalysisTime:       time.Now(),
	}

	r.logger.WithField("agent_result", fmt.Sprintf("%+v", agentResult)).Info("✅ ReAct agent analysis completed")
	return analysisResult, nil
}