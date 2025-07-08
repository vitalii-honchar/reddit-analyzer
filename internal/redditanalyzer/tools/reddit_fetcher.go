package tools

import (
	"context"
	"fmt"
	"reddit-analyzer/internal/redditanalyzer/domain"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

// TopSubredditPostsTool fetches posts from multiple subreddits
type TopSubredditPostsTool struct {
	client *reddit.Client
}

// TopSubredditPostsResult represents the result of fetching Reddit posts
type TopSubredditPostsResult struct {
	llm.BaseLLMToolResult
	Posts []domain.RedditPost `json:"posts"`
	Count int                 `json:"count"`
}

// NewTopSubredditPostsTool creates a new TopSubredditPostsTool
func NewTopSubredditPostsTool() *TopSubredditPostsTool {
	return &TopSubredditPostsTool{
		client: reddit.DefaultClient(),
	}
}

// CreateLLMTool creates the LLM tool for fetching Reddit posts
func (t *TopSubredditPostsTool) CreateLLMTool() llm.LLMTool {
	return llm.NewLLMTool(
		llm.WithLLMToolName("fetch_reddit_posts"),
		llm.WithLLMToolDescription("Fetches recent posts from multiple subreddits (last 7 days)"),
		llm.WithLLMToolParametersSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"subreddits": map[string]any{
					"type":        "array",
					"items":       map[string]any{"type": "string"},
					"description": "List of subreddit names to fetch posts from",
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Maximum number of posts to fetch per subreddit (default: 25)",
					"default":     25,
				},
				"timeframe": map[string]any{
					"type":        "string",
					"description": "Time period for posts (week, day, month)",
					"default":     "week",
				},
			},
			"required": []string{"subreddits"},
		}),
		llm.WithLLMToolCall(t.fetchPosts),
	)
}

// fetchPosts performs the actual Reddit posts fetching
func (t *TopSubredditPostsTool) fetchPosts(id string, args map[string]any) (TopSubredditPostsResult, error) {
	subreddits, ok := args["subreddits"].([]interface{})
	if !ok {
		return TopSubredditPostsResult{}, fmt.Errorf("subreddits must be an array")
	}

	limit := 25
	if l, exists := args["limit"]; exists {
		if limitFloat, ok := l.(float64); ok {
			limit = int(limitFloat)
		}
	}

	timeframe := "week"
	if t, exists := args["timeframe"]; exists {
		if timeStr, ok := t.(string); ok {
			timeframe = timeStr
		}
	}

	var allPosts []domain.RedditPost
	cutoffTime := time.Now().AddDate(0, 0, -7) // 7 days ago

	for _, subredditInterface := range subreddits {
		subredditName, ok := subredditInterface.(string)
		if !ok {
			continue
		}

		posts, err := t.fetchSubredditPosts(subredditName, limit, timeframe)
		if err != nil {
			// Log error but continue with other subreddits
			fmt.Printf("Error fetching posts from r/%s: %v\n", subredditName, err)
			continue
		}

		// Filter posts by age (last 7 days)
		for _, post := range posts {
			if post.Created.After(cutoffTime) {
				allPosts = append(allPosts, post)
			}
		}
	}

	return TopSubredditPostsResult{
		BaseLLMToolResult: llm.BaseLLMToolResult{ID: id},
		Posts:             allPosts,
		Count:             len(allPosts),
	}, nil
}

// fetchSubredditPosts fetches posts from a single subreddit
func (t *TopSubredditPostsTool) fetchSubredditPosts(subreddit string, limit int, timeframe string) ([]domain.RedditPost, error) {
	// Clean subreddit name - remove "r/" prefix if present
	cleanSubreddit := strings.TrimPrefix(subreddit, "r/")
	cleanSubreddit = strings.TrimPrefix(cleanSubreddit, "/")
	opts := &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: limit,
		},
		Time: timeframe,
	}

	posts, _, err := t.client.Subreddit.TopPosts(context.Background(), cleanSubreddit, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts from r/%s: %w", cleanSubreddit, err)
	}

	var result []domain.RedditPost
	for _, post := range posts {
		redditPost := domain.RedditPost{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Body,
			Subreddit:   subreddit,
			Author:      post.Author,
			Score:       post.Score,
			UpvoteRatio: float64(post.UpvoteRatio),
			NumComments: post.NumberOfComments,
			Created:     post.Created.Time,
			URL:         post.URL,
			Permalink:   fmt.Sprintf("https://reddit.com%s", post.Permalink),
			IsSelfPost:  post.IsSelfPost,
			Flair:       "", // Flair field needs to be properly mapped
		}
		result = append(result, redditPost)
	}

	return result, nil
}