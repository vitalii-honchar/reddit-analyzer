package tools

import (
	"reddit-analyzer/internal/redditanalyzer/domain"
	"sort"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

// FilterPostsTool filters Reddit posts by engagement metrics and criteria
type FilterPostsTool struct{}

// FilterPostsParams represents the parameters for filtering posts
type FilterPostsParams struct {
	Posts           []domain.RedditPost `json:"posts"`
	MinScore        int                 `json:"min_score,omitempty"`
	MinComments     int                 `json:"min_comments,omitempty"`
	RequireSelfPost bool                `json:"require_self_post,omitempty"`
}

// FilterPostsResult represents the result of filtering posts
type FilterPostsResult struct {
	llm.BaseLLMToolResult
	FilteredPosts []domain.RedditPost   `json:"filtered_posts"`
	OriginalCount int                   `json:"original_count"`
	FilteredCount int                   `json:"filtered_count"`
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
		llm.WithLLMToolParametersSchema[FilterPostsParams](),
		llm.WithLLMToolCall(f.filterPosts),
	)
}

// filterPosts performs the actual post filtering
func (f *FilterPostsTool) filterPosts(id string, params FilterPostsParams) (FilterPostsResult, error) {
	posts := params.Posts

	// Parse filter criteria with defaults
	criteria := domain.FilterCriteria{
		MinScore:        10,
		MinComments:     5,
		RequireSelfPost: false,
	}

	if params.MinScore > 0 {
		criteria.MinScore = params.MinScore
	}
	if params.MinComments > 0 {
		criteria.MinComments = params.MinComments
	}
	criteria.RequireSelfPost = params.RequireSelfPost

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

