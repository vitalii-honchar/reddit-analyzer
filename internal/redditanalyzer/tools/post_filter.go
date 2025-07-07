package tools

import (
	"fmt"
	"reddit-analyzer/internal/redditanalyzer/domain"
	"sort"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

// FilterPostsTool filters Reddit posts by engagement metrics and criteria
type FilterPostsTool struct{}

// FilterPostsResult represents the result of filtering posts
type FilterPostsResult struct {
	llm.BaseLLMToolResult
	FilteredPosts []domain.RedditPost `json:"filtered_posts"`
	OriginalCount int                 `json:"original_count"`
	FilteredCount int                 `json:"filtered_count"`
	Criteria      domain.FilterCriteria `json:"criteria"`
}

// NewFilterPostsTool creates a new FilterPostsTool
func NewFilterPostsTool() *FilterPostsTool {
	return &FilterPostsTool{}
}

// CreateLLMTool creates the LLM tool for filtering posts
func (f *FilterPostsTool) CreateLLMTool() llm.LLMTool {
	return llm.NewLLMTool(
		llm.WithLLMToolName("filter_posts"),
		llm.WithLLMToolDescription("Filters Reddit posts by engagement metrics and criteria"),
		llm.WithLLMToolParametersSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"posts": map[string]any{
					"type":        "array",
					"description": "Array of Reddit posts to filter",
				},
				"min_score": map[string]any{
					"type":        "number",
					"description": "Minimum upvote score (default: 10)",
					"default":     10,
				},
				"min_comments": map[string]any{
					"type":        "number",
					"description": "Minimum number of comments (default: 5)",
					"default":     5,
				},
				"require_self_post": map[string]any{
					"type":        "boolean",
					"description": "Only include self posts (text posts, not links) (default: false)",
					"default":     false,
				},
			},
			"required": []string{"posts"},
		}),
		llm.WithLLMToolCall(f.filterPosts),
	)
}

// filterPosts performs the actual post filtering
func (f *FilterPostsTool) filterPosts(id string, args map[string]any) (FilterPostsResult, error) {
	postsInterface, ok := args["posts"].([]interface{})
	if !ok {
		return FilterPostsResult{}, fmt.Errorf("posts must be an array")
	}

	// Parse filter criteria
	criteria := domain.FilterCriteria{
		MinScore:        10,
		MinComments:     5,
		RequireSelfPost: false,
	}

	if minScore, exists := args["min_score"]; exists {
		if score, ok := minScore.(float64); ok {
			criteria.MinScore = int(score)
		}
	}

	if minComments, exists := args["min_comments"]; exists {
		if comments, ok := minComments.(float64); ok {
			criteria.MinComments = int(comments)
		}
	}

	if requireSelf, exists := args["require_self_post"]; exists {
		if require, ok := requireSelf.(bool); ok {
			criteria.RequireSelfPost = require
		}
	}

	// Convert interface{} posts to domain.RedditPost
	var posts []domain.RedditPost
	for _, postInterface := range postsInterface {
		if postData, ok := postInterface.(map[string]interface{}); ok {
			post, err := f.convertToRedditPost(postData)
			if err != nil {
				continue // Skip invalid posts
			}
			posts = append(posts, post)
		}
	}

	originalCount := len(posts)
	filteredPosts := f.applyFilters(posts, criteria)

	// Sort by engagement score (score + comments)
	sort.Slice(filteredPosts, func(i, j int) bool {
		scoreI := filteredPosts[i].Score + filteredPosts[i].NumComments
		scoreJ := filteredPosts[j].Score + filteredPosts[j].NumComments
		return scoreI > scoreJ
	})

	return FilterPostsResult{
		BaseLLMToolResult: llm.BaseLLMToolResult{ID: id},
		FilteredPosts:     filteredPosts,
		OriginalCount:     originalCount,
		FilteredCount:     len(filteredPosts),
		Criteria:          criteria,
	}, nil
}

// applyFilters applies the filtering criteria to posts
func (f *FilterPostsTool) applyFilters(posts []domain.RedditPost, criteria domain.FilterCriteria) []domain.RedditPost {
	var filtered []domain.RedditPost

	for _, post := range posts {
		// Apply minimum score filter
		if post.Score < criteria.MinScore {
			continue
		}

		// Apply minimum comments filter
		if post.NumComments < criteria.MinComments {
			continue
		}

		// Apply self-post filter if required
		if criteria.RequireSelfPost && !post.IsSelfPost {
			continue
		}

		// Filter out posts that are too short (likely not substantive)
		if len(post.Title) < 10 {
			continue
		}

		// Filter out deleted/removed posts
		if post.Author == "[deleted]" || post.Content == "[removed]" {
			continue
		}

		filtered = append(filtered, post)
	}

	return filtered
}

// convertToRedditPost converts a map to a RedditPost struct
func (f *FilterPostsTool) convertToRedditPost(data map[string]interface{}) (domain.RedditPost, error) {
	post := domain.RedditPost{}

	if id, ok := data["id"].(string); ok {
		post.ID = id
	}
	if title, ok := data["title"].(string); ok {
		post.Title = title
	}
	if content, ok := data["content"].(string); ok {
		post.Content = content
	}
	if subreddit, ok := data["subreddit"].(string); ok {
		post.Subreddit = subreddit
	}
	if author, ok := data["author"].(string); ok {
		post.Author = author
	}
	if score, ok := data["score"].(float64); ok {
		post.Score = int(score)
	}
	if ratio, ok := data["upvote_ratio"].(float64); ok {
		post.UpvoteRatio = ratio
	}
	if comments, ok := data["num_comments"].(float64); ok {
		post.NumComments = int(comments)
	}
	if url, ok := data["url"].(string); ok {
		post.URL = url
	}
	if permalink, ok := data["permalink"].(string); ok {
		post.Permalink = permalink
	}
	if isSelf, ok := data["is_self_post"].(bool); ok {
		post.IsSelfPost = isSelf
	}
	if flair, ok := data["flair"].(string); ok {
		post.Flair = flair
	}

	return post, nil
}