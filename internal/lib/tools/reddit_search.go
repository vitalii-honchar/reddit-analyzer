package tools

import (
	"context"
	"fmt"
	"reddit-analyzer/internal/lib/domain"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
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

const (
	defaultSubredditPostSearchLimit   = 25
	defaultSubredditPostSearchTimeout = 60 * time.Second
)

type (
	SubredditPostSearchParams struct {
		Subreddit       string                     `json:"subreddit" jsonschema_description:"Subreddit name (without r/ prefix) to search in"`
		Filter          FilterCriteria             `json:"filter" jsonschema_description:"Criteria to filter posts for relevance and quality"`
		Type            SubredditPostType          `json:"type" json_schema:"enum=top,enum=hot,enum=new,enum=controversial" jsonschema_description:"Type of posts to fetch (required)"`
		Timeframe       SubredditPostTypeTimeframe `json:"timeframe" json_schema:"enum=hour,enum=day,enum=week,enum=month,enum=year,enum=all" jsonschema_description:"Timeframe for posts to fetch - required for top and controversial post types"`
		Limit           int                        `json:"limit,omitempty" jsonschema_description:"Maximum number of posts to fetch (default: 25, used to control API usage)"`
		PaginationToken string                     `json:"pagination_token,omitempty" jsonschema_description:"Token for paginating results if more posts are available"`
		Timeout         time.Duration              `json:"timeout,omitempty" jsonschema_description:"Timeout for the search operation, default is 60 seconds"`
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

	redditPostSearchTool struct {
		client *reddit.Client
	}
)

func NewRedditPostSearchTool() llm.LLMTool {
	tool := &redditPostSearchTool{
		client: reddit.DefaultClient(),
	}

	return llm.NewLLMTool(
		llm.WithLLMToolName("reddit_post_search"),
		llm.WithLLMToolDescription("Searches Reddit for posts matching specific criteria"),
		llm.WithLLMToolParametersSchema[SubredditPostSearchParams](),
		llm.WithLLMToolCall(tool.search),
	)
}

func (r *redditPostSearchTool) search(callID string, params SubredditPostSearchParams) (SubredditPostSearchResult, error) {

	if err := r.validateParams(&params); err != nil {
		return SubredditPostSearchResult{}, err
	}

	posts, err := r.searchPosts(params)
	if err != nil {
		return SubredditPostSearchResult{}, err
	}

	filteredPosts, err := r.filterPosts(posts, params.Filter)
	if err != nil {
		return SubredditPostSearchResult{}, err
	}

	result := SubredditPostSearchResult{
		BaseLLMToolResult: llm.BaseLLMToolResult{
			ID: callID,
		},
		Posts:         filteredPosts,
		CountAnalyzed: len(posts),
	}

	// Set pagination token if we have posts
	if len(posts) > 0 {
		result.PaginationToken = posts[len(posts)-1].FullID
	}

	return result, nil
}

func (r *redditPostSearchTool) searchPosts(params SubredditPostSearchParams) ([]domain.RedditPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), params.Timeout)
	defer cancel()

	// Clean subreddit name - remove "r/" prefix if present
	cleanSubreddit := strings.TrimPrefix(params.Subreddit, "r/")
	cleanSubreddit = strings.TrimPrefix(cleanSubreddit, "/")

	// Set up list options
	opts := &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: params.Limit,
		},
	}

	// Add pagination token if provided
	if params.PaginationToken != "" {
		opts.ListOptions.After = params.PaginationToken
	}

	// Add timeframe for top and controversial posts
	if params.Type == TopSubredditPostType || params.Type == ControversialSubredditPostType {
		opts.Time = string(params.Timeframe)
	}

	// Fetch posts based on type
	var posts []*reddit.Post
	var err error

	switch params.Type {
	case TopSubredditPostType:
		posts, _, err = r.client.Subreddit.TopPosts(ctx, cleanSubreddit, opts)
	case HotSubredditPostType:
		posts, _, err = r.client.Subreddit.HotPosts(ctx, cleanSubreddit, &opts.ListOptions)
	case NewSubredditPostType:
		posts, _, err = r.client.Subreddit.NewPosts(ctx, cleanSubreddit, &opts.ListOptions)
	case ControversialSubredditPostType:
		posts, _, err = r.client.Subreddit.ControversialPosts(ctx, cleanSubreddit, opts)
	default:
		return nil, fmt.Errorf("%w: invalid post type %s", llm.ErrInvalidArguments, params.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts from r/%s: %w", cleanSubreddit, err)
	}

	// Convert Reddit posts to domain posts
	var result []domain.RedditPost
	for _, post := range posts {
		redditPost := domain.RedditPost{
			ID:                    post.ID,
			FullID:                post.FullID,
			Created:               post.Created.Time,
			Edited:                post.Edited.Time,
			Permalink:             fmt.Sprintf("https://reddit.com%s", post.Permalink),
			URL:                   post.URL,
			Title:                 post.Title,
			Body:                  post.Body,
			Score:                 post.Score,
			UpvoteRatio:           post.UpvoteRatio,
			NumberOfComments:      post.NumberOfComments,
			SubredditName:         post.SubredditName,
			SubredditNamePrefixed: post.SubredditNamePrefixed,
			SubredditID:           post.SubredditID,
			SubredditSubscribers:  post.SubredditSubscribers,
			Author:                post.Author,
			AuthorID:              post.AuthorID,
			Spoiler:               post.Spoiler,
			Locked:                post.Locked,
			NSFW:                  post.NSFW,
			IsSelfPost:            post.IsSelfPost,
			Saved:                 post.Saved,
			Stickied:              post.Stickied,
		}
		result = append(result, redditPost)
	}

	return result, nil
}

func (r *redditPostSearchTool) filterPosts(posts []domain.RedditPost, criteria FilterCriteria) ([]domain.RedditPost, error) {
	var filtered []domain.RedditPost
	cutoffTime := time.Now().Add(-criteria.MaxAge)

	for _, post := range posts {
		// Apply minimum score filter
		if post.Score < criteria.MinScore {
			continue
		}

		// Apply minimum comments filter
		if post.NumberOfComments < criteria.MinComments {
			continue
		}

		// Apply max age filter
		if criteria.MaxAge > 0 && post.Created.Before(cutoffTime) {
			continue
		}

		// Apply stickied filter
		if criteria.ExcludeStickied && post.Stickied {
			continue
		}

		// Apply self post filter
		if criteria.RequireSelfPost && !post.IsSelfPost {
			continue
		}

		filtered = append(filtered, post)
	}

	return filtered, nil
}

func (r *redditPostSearchTool) validateParams(params *SubredditPostSearchParams) error {
	// Validate subreddit name
	if params.Subreddit == "" {
		return fmt.Errorf("%w: subreddit name cannot be empty", llm.ErrInvalidArguments)
	}

	if strings.HasPrefix(params.Subreddit, "r/") {
		return fmt.Errorf("%w: subreddit %s should not have r/ prefix", llm.ErrInvalidArguments, params.Subreddit)
	}

	// Set default values
	if params.Limit <= 0 {
		params.Limit = defaultSubredditPostSearchLimit
	}

	if params.Timeout <= 0 {
		params.Timeout = defaultSubredditPostSearchTimeout
	}

	// Validate post type is specified
	if params.Type == "" {
		return fmt.Errorf("%w: post type is required", llm.ErrInvalidArguments)
	}

	// Validate timeframe is specified for top and controversial posts
	if (params.Type == TopSubredditPostType || params.Type == ControversialSubredditPostType) && params.Timeframe == "" {
		return fmt.Errorf("%w: timeframe is required for top and controversial post types", llm.ErrInvalidArguments)
	}

	// Set default filter criteria if not provided
	if params.Filter.MinScore == 0 {
		params.Filter.MinScore = 1
	}
	if params.Filter.MinComments == 0 {
		params.Filter.MinComments = 1
	}
	if params.Filter.MaxAge == 0 {
		params.Filter.MaxAge = 7 * 24 * time.Hour // 7 days default
	}

	return nil
}
