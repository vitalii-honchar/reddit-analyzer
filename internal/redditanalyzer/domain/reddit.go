package domain

import (
	"time"
)

// RedditPost represents a Reddit post with relevant fields for analysis
type RedditPost struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Subreddit     string    `json:"subreddit"`
	Author        string    `json:"author"`
	Score         int       `json:"score"`
	UpvoteRatio   float64   `json:"upvote_ratio"`
	NumComments   int       `json:"num_comments"`
	Created       time.Time `json:"created"`
	URL           string    `json:"url"`
	Permalink     string    `json:"permalink"`
	IsSelfPost    bool      `json:"is_self_post"`
	Flair         string    `json:"flair,omitempty"`
}

// SubredditSelection represents AI-selected subreddits for analysis
type SubredditSelection struct {
	Subreddits []string `json:"subreddits"`
	Reasoning  string   `json:"reasoning"`
}

// FilterCriteria defines criteria for filtering posts
type FilterCriteria struct {
	MinScore         int           `json:"min_score"`
	MinComments      int           `json:"min_comments"`
	MaxAge           time.Duration `json:"max_age"`
	ExcludeStickied  bool          `json:"exclude_stickied"`
	RequireSelfPost  bool          `json:"require_self_post"`
}