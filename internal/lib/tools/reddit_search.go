package tools

import (
	"reddit-analyzer/internal/lib/domain"
	"time"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

type SubredditPostType string

const (
	TopSubredditPostType           = SubredditPostType("top")
	HotSubredditPostType           = SubredditPostType("hot")
	NewSubredditPostType           = SubredditPostType("new")
	ControversialSubredditPostType = SubredditPostType("controversial")
)

type SubredditPostTypeTimeframe string

const (
	HourSubredditPostTypeTimeframe    = SubredditPostTypeTimeframe("hour")
	DaySubredditPostTypeTimeframe     = SubredditPostTypeTimeframe("day")
	WeekSubredditPostTypeTimeframe    = SubredditPostTypeTimeframe("week")
	MonthSubredditPostTypeTimeframe   = SubredditPostTypeTimeframe("month")
	YearSubredditPostTypeTimeframe    = SubredditPostTypeTimeframe("year")
	AllTimeSubredditPostTypeTimeframe = SubredditPostTypeTimeframe("all")
)

const defaultSubredditPostLimit = 25

type (
	SubredditPostSearchParams struct {
		Subreddit       string                     `json:"subreddit" jsonschema_description:"Subreddit name (without r/ prefix) to search in"`
		Filter          FilterCriteria             `json:"filter,omitempty" jsonschema_description:"Criteria to filter posts for relevance and quality"`
		Type            SubredditPostType          `json:"type,omitempty" json_schema:"enum=top,enum=hot,enum=new,enum=controversial" jsonschema_description:"Type of posts to fetch"`
		Timeframe       SubredditPostTypeTimeframe `json:"timeframe,omitempty" json_schema:"enum=hour,enum=day,enum=week,enum=month,enum=year,enum=all" jsonschema_description:"Timeframe for posts to fetch"`
		Limit           int                        `json:"limit,omitempty" jsonschema_description:"Maximum number of posts to fetch (default: 25, used to control API usage)"`
		PaginationToken string                     `json:"pagination_token,omitempty" jsonschema_description:"Token for paginating results if more posts are available"`
	}

	FilterCriteria struct {
		MinScore        int           `json:"min_score" jsonschema_description:"Minimum upvote score required for posts to be considered"`
		MinComments     int           `json:"min_comments" jsonschema_description:"Minimum number of comments required to indicate community discussion"`
		MaxAge          time.Duration `json:"max_age" jsonschema_description:"Maximum age of posts to include in analysis (e.g. 7 days)"`
		ExcludeStickied bool          `json:"exclude_stickied" jsonschema_description:"Whether to exclude moderator-pinned posts from analysis"`
		RequireSelfPost bool          `json:"require_self_post" jsonschema_description:"Whether to only include text posts (not links) which often contain more detailed problems"`
	}

	SubredditPostSearchResult struct {
		llm.BaseLLMToolResult
		Posts           []domain.RedditPost `json:"posts" jsonschema_description:"Array of Reddit posts matching the search criteria"`
		CountAnalyzed   int                 `json:"count_analyzed" jsonschema_description:"Total number of posts analyzed from the search"`
		PaginationToken string              `json:"pagination_token,omitempty" jsonschema_description:"Token for paginating results if more posts are available"`
	}
)

func NewRedditPostSearchTool() llm.LLMTool {
	return llm.NewLLMTool(
		llm.WithLLMToolName("reddit_post_search"),
		llm.WithLLMToolDescription("Searches Reddit for posts matching specific criteria"),
		llm.WithLLMToolParametersSchema[SubredditPostSearchParams](),
		llm.WithLLMToolCall[SubredditPostSearchParams, SubredditPostSearchResult](redditPostSearchTool()),
	)
}

func redditPostSearchTool() func(string, SubredditPostSearchParams) (SubredditPostSearchResult, error) {
	return func(callID string, params SubredditPostSearchParams) (SubredditPostSearchResult, error) {
		// Implement the logic to search Reddit posts based on the provided parameters
		// and populate the result with the found posts.
		return SubredditPostSearchResult{}, nil
	}
}
