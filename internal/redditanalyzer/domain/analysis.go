package domain

import (
	"time"
)

// OpportunityAnalysis represents the LLM analysis of a Reddit post
type OpportunityAnalysis struct {
	PostID          string  `json:"post_id"`
	Score           int     `json:"score"`           // 1-5 rating for indie hacker opportunity
	ProblemSummary  string  `json:"problem_summary"` // Brief description of the problem
	MarketAnalysis  string  `json:"market_analysis"` // Analysis of market gap/opportunity
	TargetAudience  string  `json:"target_audience"` // Who has this problem
	CompetitionNote string  `json:"competition_note"` // Notes about existing solutions
	Reasoning       string  `json:"reasoning"`       // Why this score was given
}

// AnalysisResult represents the final result of the Reddit analysis
type AnalysisResult struct {
	ProjectDirection string                `json:"project_direction"`
	SelectedSubreddits []string           `json:"selected_subreddits"`
	PostsAnalyzed    int                  `json:"posts_analyzed"`
	PostsFiltered    int                  `json:"posts_filtered"`
	Opportunities    []RankedOpportunity  `json:"opportunities"`
	AnalysisTime     time.Time            `json:"analysis_time"`
	ExecutionTime    time.Duration        `json:"execution_time"`
}

// RankedOpportunity combines a Reddit post with its opportunity analysis
type RankedOpportunity struct {
	Post     RedditPost          `json:"post"`
	Analysis OpportunityAnalysis `json:"analysis"`
	Rank     int                 `json:"rank"` // Overall ranking (1 = best opportunity)
}

// AgentInput represents the input provided to the master agent
type AgentInput struct {
	ProjectDirection string `json:"project_direction"`
	MaxPosts         int    `json:"max_posts,omitempty"`
	MaxAge           string `json:"max_age,omitempty"` // e.g., "7d", "1w"
}