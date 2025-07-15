package domain

import "time"

type RedditPost struct {
	ID                    string    `json:"id,omitempty"`
	FullID                string    `json:"name,omitempty"`
	Created               time.Time `json:"created_utc,omitempty"`
	Edited                time.Time `json:"edited,omitempty"`
	Permalink             string    `json:"permalink,omitempty"`
	URL                   string    `json:"url,omitempty"`
	Title                 string    `json:"title,omitempty"`
	Body                  string    `json:"selftext,omitempty"`
	Score                 int       `json:"score"`
	UpvoteRatio           float32   `json:"upvote_ratio"`
	NumberOfComments      int       `json:"num_comments"`
	SubredditName         string    `json:"subreddit,omitempty"`
	SubredditNamePrefixed string    `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string    `json:"subreddit_id,omitempty"`
	SubredditSubscribers  int       `json:"subreddit_subscribers"`
	Author                string    `json:"author,omitempty"`
	AuthorID              string    `json:"author_fullname,omitempty"`
	Spoiler               bool      `json:"spoiler"`
	Locked                bool      `json:"locked"`
	NSFW                  bool      `json:"over_18"`
	IsSelfPost            bool      `json:"is_self"`
	Saved                 bool      `json:"saved"`
	Stickied              bool      `json:"stickied"`
}

type FilterCriteria struct {
	MinScore        int           `json:"min_score" jsonschema_description:"Minimum upvote score required for posts to be considered (default: 10)"`
	MinComments     int           `json:"min_comments" jsonschema_description:"Minimum number of comments required to indicate community discussion (default: 5)"`
	MaxAge          time.Duration `json:"max_age" jsonschema_description:"Maximum age of posts to include in analysis (e.g. 7 days)"`
	ExcludeStickied bool          `json:"exclude_stickied" jsonschema_description:"Whether to exclude moderator-pinned posts from analysis"`
	RequireSelfPost bool          `json:"require_self_post" jsonschema_description:"Whether to only include text posts (not links) which often contain more detailed problems"`
}
