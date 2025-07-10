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

// TopSubredditPostsParams represents the parameters for fetching Reddit posts
type TopSubredditPostsParams struct {
	Subreddits []string `json:"subreddits"`
	Limit      int      `json:"limit,omitempty"`
	Timeframe  string   `json:"timeframe,omitempty"`
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
		llm.WithLLMToolParametersSchema[TopSubredditPostsParams](),
		llm.WithLLMToolCall(t.fetchPosts),
	)
}

// fetchPosts performs the actual Reddit posts fetching
func (t *TopSubredditPostsTool) fetchPosts(id string, params TopSubredditPostsParams) (TopSubredditPostsResult, error) {
	subreddits := params.Subreddits
	
	limit := 25
	if params.Limit > 0 {
		limit = params.Limit
	}

	timeframe := "week"
	if params.Timeframe != "" {
		timeframe = params.Timeframe
	}

	var allPosts []domain.RedditPost
	cutoffTime := time.Now().AddDate(0, 0, -7) // 7 days ago

	for _, subredditName := range subreddits {

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
			Flair:       nil, // Flair field is optional, set to nil
		}
		result = append(result, redditPost)
	}

	return result, nil
}