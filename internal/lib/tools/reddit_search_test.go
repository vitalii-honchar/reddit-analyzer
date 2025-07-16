package tools

import (
	"reddit-analyzer/internal/lib/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func TestRedditPostSearchTool_Integration(t *testing.T) {
	t.Parallel()

	// Create the tool directly to access its methods
	tool := &redditPostSearchTool{
		client: reddit.DefaultClient(),
	}

	t.Run("search golang subreddit hot posts", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Limit:     5,
			Filter: FilterCriteria{
				MinScore:        1,
				MinComments:     1,
				MaxAgeDays:      7 * 24 * time.Hour,
				ExcludeStickied: false,
				RequireSelfPost: false,
			},
		}

		result, err := tool.search("test-call-id", params)
		require.NoError(t, err)
		assert.Equal(t, "test-call-id", result.ID)
		assert.LessOrEqual(t, len(result.Posts), 5)
		assert.GreaterOrEqual(t, result.CountAnalyzed, len(result.Posts))

		// Verify posts have required fields
		for _, post := range result.Posts {
			assert.NotEmpty(t, post.ID)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.SubredditName)
			assert.GreaterOrEqual(t, post.Score, 1)
			assert.GreaterOrEqual(t, post.NumberOfComments, 1)
		}
	})

	t.Run("search golang subreddit top posts with timeframe", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      TopSubredditPostType,
			Timeframe: WeekSubredditPostTypeTimeframe,
			Limit:     3,
			Filter: FilterCriteria{
				MinScore:        5,
				MinComments:     2,
				MaxAgeDays:      7 * 24 * time.Hour,
				ExcludeStickied: true,
				RequireSelfPost: false,
			},
		}

		result, err := tool.search("test-call-id-2", params)
		require.NoError(t, err)
		assert.Equal(t, "test-call-id-2", result.ID)
		assert.LessOrEqual(t, len(result.Posts), 3)

		// Verify filtering worked
		for _, post := range result.Posts {
			assert.GreaterOrEqual(t, post.Score, 5)
			assert.GreaterOrEqual(t, post.NumberOfComments, 2)
			assert.False(t, post.Stickied, "Should exclude stickied posts")
		}
	})

	t.Run("search with self post requirement", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      NewSubredditPostType,
			Limit:     10,
			Filter: FilterCriteria{
				MinScore:        1,
				MinComments:     1,
				MaxAgeDays:      30 * 24 * time.Hour, // 30 days
				ExcludeStickied: true,
				RequireSelfPost: true,
			},
		}

		result, err := tool.search("test-call-id-3", params)
		require.NoError(t, err)
		assert.Equal(t, "test-call-id-3", result.ID)

		// Verify all returned posts are self posts
		for _, post := range result.Posts {
			assert.True(t, post.IsSelfPost, "Should only return self posts")
		}
	})

	t.Run("pagination with token", func(t *testing.T) {
		t.Parallel()

		// First request to get initial posts and pagination token
		firstParams := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Limit:     3, // Small limit to ensure pagination
			Filter: FilterCriteria{
				MinScore:        1,
				MinComments:     1,
				MaxAgeDays:      7 * 24 * time.Hour,
				ExcludeStickied: false,
				RequireSelfPost: false,
			},
		}

		firstResult, err := tool.search("test-pagination-1", firstParams)
		require.NoError(t, err)
		assert.Equal(t, "test-pagination-1", firstResult.ID)
		assert.LessOrEqual(t, len(firstResult.Posts), 3)
		assert.NotEmpty(t, firstResult.PaginationToken, "Should return pagination token")

		// Second request using pagination token
		secondParams := SubredditPostSearchParams{
			Subreddit:       "golang",
			Type:            HotSubredditPostType,
			Limit:           3,
			PaginationToken: firstResult.PaginationToken,
			Filter: FilterCriteria{
				MinScore:        1,
				MinComments:     1,
				MaxAgeDays:      7 * 24 * time.Hour,
				ExcludeStickied: false,
				RequireSelfPost: false,
			},
		}

		secondResult, err := tool.search("test-pagination-2", secondParams)
		require.NoError(t, err)
		assert.Equal(t, "test-pagination-2", secondResult.ID)
		assert.LessOrEqual(t, len(secondResult.Posts), 3)

		// Verify we got different posts (pagination worked)
		if len(firstResult.Posts) > 0 && len(secondResult.Posts) > 0 {
			// Check that first post IDs are different
			firstPostID := firstResult.Posts[0].ID
			secondPostID := secondResult.Posts[0].ID
			assert.NotEqual(t, firstPostID, secondPostID, "Pagination should return different posts")
		}
	})
}

func TestRedditPostSearchTool_ValidateParams(t *testing.T) {
	t.Parallel()

	tool := &redditPostSearchTool{}

	t.Run("valid params with defaults", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    5,
				MinComments: 3,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.NoError(t, err)

		// Check defaults were set for optional parameters
		assert.Equal(t, defaultSubredditPostSearchLimit, params.Limit)
		assert.Equal(t, defaultSubredditPostSearchTimeout, params.Timeout)
		assert.Equal(t, HotSubredditPostType, params.Type)
		// Check required filter criteria were preserved
		assert.Equal(t, 5, params.Filter.MinScore)
		assert.Equal(t, 3, params.Filter.MinComments)
		assert.Equal(t, 7*24*time.Hour, params.Filter.MaxAgeDays)
	})

	t.Run("top posts require timeframe", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      TopSubredditPostType,
			Timeframe: WeekSubredditPostTypeTimeframe,
			Filter: FilterCriteria{
				MinScore:    10,
				MinComments: 5,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.NoError(t, err)
		assert.Equal(t, WeekSubredditPostTypeTimeframe, params.Timeframe)
	})

	t.Run("controversial posts require timeframe", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      ControversialSubredditPostType,
			Timeframe: WeekSubredditPostTypeTimeframe,
			Filter: FilterCriteria{
				MinScore:    15,
				MinComments: 8,
				MaxAgeDays:  3 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.NoError(t, err)
		assert.Equal(t, WeekSubredditPostTypeTimeframe, params.Timeframe)
	})

	t.Run("empty subreddit name fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "subreddit name cannot be empty")
	})

	t.Run("empty post type fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      "",
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "post type is required")
	})

	t.Run("top posts without timeframe fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      TopSubredditPostType,
			Timeframe: "",
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeframe is required for top and controversial post types")
	})

	t.Run("controversial posts without timeframe fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      ControversialSubredditPostType,
			Timeframe: "",
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeframe is required for top and controversial post types")
	})

	t.Run("invalid post type fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      SubredditPostType("invalid"),
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid post type invalid")
	})

	t.Run("invalid timeframe fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      TopSubredditPostType,
			Timeframe: SubredditPostTypeTimeframe("invalid"),
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid timeframe invalid")
	})

	t.Run("negative min_score fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    -1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "min_score must be >= 0")
	})

	t.Run("negative min_comments fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: -1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "min_comments must be >= 0")
	})

	t.Run("zero max_age fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  0,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "max_age must be > 0")
	})

	t.Run("negative max_age fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  -24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "max_age must be > 0")
	})

	t.Run("subreddit with r/ prefix fails", func(t *testing.T) {
		t.Parallel()

		params := SubredditPostSearchParams{
			Subreddit: "r/golang",
			Type:      HotSubredditPostType,
			Filter: FilterCriteria{
				MinScore:    1,
				MinComments: 1,
				MaxAgeDays:  7 * 24 * time.Hour,
			},
		}

		err := tool.validateParams(&params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "should not have r/ prefix")
	})
}

func TestRedditPostSearchTool_FilterPosts(t *testing.T) {
	t.Parallel()

	tool := &redditPostSearchTool{}

	// Create test posts
	now := time.Now()
	posts := []domain.RedditPost{
		{
			ID:               "post1",
			Title:            "High score post",
			Score:            100,
			NumberOfComments: 50,
			Created:          now.Add(-1 * time.Hour),
			IsSelfPost:       true,
			Stickied:         false,
		},
		{
			ID:               "post2",
			Title:            "Low score post",
			Score:            2,
			NumberOfComments: 1,
			Created:          now.Add(-2 * time.Hour),
			IsSelfPost:       false,
			Stickied:         false,
		},
		{
			ID:               "post3",
			Title:            "Old post",
			Score:            50,
			NumberOfComments: 25,
			Created:          now.Add(-10 * 24 * time.Hour), // 10 days old
			IsSelfPost:       true,
			Stickied:         false,
		},
		{
			ID:               "post4",
			Title:            "Stickied post",
			Score:            75,
			NumberOfComments: 30,
			Created:          now.Add(-3 * time.Hour),
			IsSelfPost:       true,
			Stickied:         true,
		},
		{
			ID:               "post5",
			Title:            "Link post",
			Score:            60,
			NumberOfComments: 20,
			Created:          now.Add(-4 * time.Hour),
			IsSelfPost:       false,
			Stickied:         false,
		},
	}

	t.Run("filter by minimum score and comments", func(t *testing.T) {
		t.Parallel()

		criteria := FilterCriteria{
			MinScore:    50,
			MinComments: 20,
			MaxAgeDays:  7 * 24 * time.Hour,
		}

		filtered, err := tool.filterPosts(posts, criteria)
		require.NoError(t, err)

		// Should include post1, post4, and post5 (post3 is too old)
		assert.Len(t, filtered, 3)
		assert.Contains(t, getPostIDs(filtered), "post1")
		assert.Contains(t, getPostIDs(filtered), "post4")
		assert.Contains(t, getPostIDs(filtered), "post5")
	})

	t.Run("exclude stickied posts", func(t *testing.T) {
		t.Parallel()

		criteria := FilterCriteria{
			MinScore:        10,
			MinComments:     5,
			MaxAgeDays:      7 * 24 * time.Hour,
			ExcludeStickied: true,
		}

		filtered, err := tool.filterPosts(posts, criteria)
		require.NoError(t, err)

		// Should not include post4 (stickied)
		postIDs := getPostIDs(filtered)
		assert.NotContains(t, postIDs, "post4")
		for _, post := range filtered {
			assert.False(t, post.Stickied)
		}
	})

	t.Run("require self posts only", func(t *testing.T) {
		t.Parallel()

		criteria := FilterCriteria{
			MinScore:        10,
			MinComments:     5,
			MaxAgeDays:      7 * 24 * time.Hour,
			RequireSelfPost: true,
		}

		filtered, err := tool.filterPosts(posts, criteria)
		require.NoError(t, err)

		// Should not include post2 and post5 (not self posts)
		postIDs := getPostIDs(filtered)
		assert.NotContains(t, postIDs, "post2")
		assert.NotContains(t, postIDs, "post5")
		for _, post := range filtered {
			assert.True(t, post.IsSelfPost)
		}
	})

	t.Run("filter by age", func(t *testing.T) {
		t.Parallel()

		criteria := FilterCriteria{
			MinScore:    10,
			MinComments: 5,
			MaxAgeDays:  2 * 24 * time.Hour, // 2 days
		}

		filtered, err := tool.filterPosts(posts, criteria)
		require.NoError(t, err)

		// Should not include post3 (too old)
		postIDs := getPostIDs(filtered)
		assert.NotContains(t, postIDs, "post3")
		for _, post := range filtered {
			assert.True(t, post.Created.After(now.Add(-2*24*time.Hour)))
		}
	})
}

// Helper function to extract post IDs from a slice of posts
func getPostIDs(posts []domain.RedditPost) []string {
	ids := make([]string, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	return ids
}
