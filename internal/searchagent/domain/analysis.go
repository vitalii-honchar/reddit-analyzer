package domain

import (
	"time"
)

// OpportunityAnalysis represents the LLM analysis of a Reddit post
type OpportunityAnalysis struct {
	PostID          string  `json:"post_id" jsonschema_description:"Reddit post ID that this analysis corresponds to"`
	Score           int     `json:"score" jsonschema_description:"Opportunity rating from 1-5 where 5=exceptional indie hacker opportunity, 1=poor opportunity"`
	ProblemSummary  string  `json:"problem_summary" jsonschema_description:"Concise summary of the core problem or pain point mentioned in the post"`
	MarketAnalysis  string  `json:"market_analysis" jsonschema_description:"Analysis of the market gap, opportunity size, and why this could be viable for indie hackers"`
	TargetAudience  string  `json:"target_audience" jsonschema_description:"Specific description of who has this problem (e.g. 'small MSPs', 'remote teams', 'sysadmins')"`
	CompetitionNote string  `json:"competition_note" jsonschema_description:"Notes about existing solutions, competitors, or lack thereof in this space"`
	Reasoning       string  `json:"reasoning" jsonschema_description:"Detailed explanation of why this score was assigned, focusing on indie hacker feasibility"`
}

// AnalysisResult represents the final result of the Reddit analysis
type AnalysisResult struct {
	ProjectDirection   string              `json:"project_direction" jsonschema_description:"The original user input describing their project direction or domain of interest"`
	SelectedSubreddits []string            `json:"selected_subreddits" jsonschema_description:"List of niche subreddits that were chosen for analysis based on the project direction"`
	PostsAnalyzed      int                 `json:"posts_analyzed" jsonschema_description:"Total number of Reddit posts that were analyzed for opportunities"`
	PostsFiltered      int                 `json:"posts_filtered" jsonschema_description:"Number of posts that passed the filtering criteria (engagement, recency, etc.)"`
	Opportunities      []RankedOpportunity `json:"opportunities" jsonschema_description:"List of discovered opportunities ranked by potential, with highest scores first"`
	AnalysisTime       time.Time           `json:"analysis_time" jsonschema_description:"Timestamp when this analysis was completed"`
	ExecutionTime      time.Duration       `json:"execution_time" jsonschema_description:"Total time taken to complete the entire analysis workflow"`
}

// RankedOpportunity combines a Reddit post with its opportunity analysis
type RankedOpportunity struct {
	Post     RedditPost          `json:"post" jsonschema_description:"The original Reddit post that contains the problem or opportunity"`
	Analysis OpportunityAnalysis `json:"analysis" jsonschema_description:"Detailed LLM analysis of why this post represents an indie hacker opportunity"`
	Rank     int                 `json:"rank" jsonschema_description:"Overall ranking position where 1 = best opportunity, 2 = second best, etc."`
}

// AgentInput represents the input provided to the master agent
type AgentInput struct {
	ProjectDirection string `json:"project_direction" jsonschema_description:"User's project direction or domain of interest (e.g. 'cybersecurity project', 'productivity tools')"`
	MaxPosts         int    `json:"max_posts,omitempty" jsonschema_description:"Optional limit on maximum number of posts to analyze (default: system decides)"`
	MaxAge           string `json:"max_age,omitempty" jsonschema_description:"Optional age limit for posts (e.g. '7d', '1w') - how far back to look for posts"`
}