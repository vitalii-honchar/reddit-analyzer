package domain

import "time"

type RedditPost struct {
	ID                    string    `json:"id"`
	FullID                string    `json:"name"`
	Created               time.Time `json:"created_utc"`
	Edited                time.Time `json:"edited"`
	Permalink             string    `json:"permalink"`
	URL                   string    `json:"url"`
	Title                 string    `json:"title"`
	Body                  string    `json:"selftext"`
	Score                 int       `json:"score,omitempty"`
	UpvoteRatio           float32   `json:"upvote_ratio,omitempty"`
	NumberOfComments      int       `json:"num_comments,omitempty"`
	SubredditName         string    `json:"subreddit,omitempty"`
	SubredditNamePrefixed string    `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string    `json:"subreddit_id"`
	SubredditSubscribers  int       `json:"subreddit_subscribers"`
	Author                string    `json:"author,omitempty"`
	AuthorID              string    `json:"author_fullname,omitempty"`
	Spoiler               bool      `json:"spoiler,omitempty"`
	Locked                bool      `json:"locked,omitempty"`
	NSFW                  bool      `json:"over_18,omitempty"`
	IsSelfPost            bool      `json:"is_self,omitempty"`
	Saved                 bool      `json:"saved,omitempty"`
	Stickied              bool      `json:"stickied,omitempty"`
}

type FilterCriteria struct {
	MinScore        int           `json:"min_score" jsonschema_description:"Minimum upvote score required for posts to be considered (default: 10)"`
	MinComments     int           `json:"min_comments" jsonschema_description:"Minimum number of comments required to indicate community discussion (default: 5)"`
	MaxAge          time.Duration `json:"max_age" jsonschema_description:"Maximum age of posts to include in analysis (e.g. 7 days)"`
	ExcludeStickied bool          `json:"exclude_stickied" jsonschema_description:"Whether to exclude moderator-pinned posts from analysis"`
	RequireSelfPost bool          `json:"require_self_post" jsonschema_description:"Whether to only include text posts (not links) which often contain more detailed problems"`
}
