package domain

type SearchResult struct {
	RedditPosts []RedditPost `json:"reddit_posts" jsonschema_description:"List of Reddit posts matching the search query"`
}
