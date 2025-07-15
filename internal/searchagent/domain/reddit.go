package domain

import (
	"time"
)

// RedditPost represents a Reddit post with relevant fields for analysis
type RedditPost struct {
	ID          string    `json:"id" jsonschema_description:"Unique Reddit post identifier (e.g. 'abc123')"`
	Title       string    `json:"title" jsonschema_description:"Post title text that often contains the main problem or question"`
	Content     string    `json:"content" jsonschema_description:"Full post content/body text where users describe their problems in detail"`
	Subreddit   string    `json:"subreddit" jsonschema_description:"Name of the subreddit where this was posted (e.g. 'sysadmin', 'cybersecurity')"`
	Author      string    `json:"author" jsonschema_description:"Reddit username of the person who posted this"`
	Score       int       `json:"score" jsonschema_description:"Net upvotes (upvotes minus downvotes) indicating community engagement"`
	UpvoteRatio float64   `json:"upvote_ratio" jsonschema_description:"Ratio of upvotes to total votes (0.0-1.0), higher means more positive reception"`
	NumComments int       `json:"num_comments" jsonschema_description:"Total number of comments, indicating discussion level and community interest"`
	Created     time.Time `json:"created" jsonschema_description:"When this post was created on Reddit"`
	URL         string    `json:"url" jsonschema_description:"Original Reddit URL for this post"`
	Permalink   string    `json:"permalink" jsonschema_description:"Permanent Reddit link to this specific post"`
	IsSelfPost  bool      `json:"is_self_post" jsonschema_description:"True if this is a text post (not a link), often containing more detailed problem descriptions"`
}

// SubredditSelection represents AI-selected subreddits for analysis
type SubredditSelection struct {
	Subreddits []string `json:"subreddits" jsonschema_description:"List of 3-5 niche subreddit names (without r/ prefix) chosen for analysis based on project direction"`
	Reasoning  string   `json:"reasoning" jsonschema_description:"Brief explanation of why these specific subreddits were selected for finding hidden opportunities"`
}

// FilterCriteria defines criteria for filtering posts
