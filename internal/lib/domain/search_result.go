package domain

type SearchResult struct {
	RedditPosts []RedditPostFinding `json:"reddit_posts" jsonschema_description:"Summary of relevant Reddit posts found during the search"`
}

type RedditPostFinding struct {
	ID        string `json:"id" jsonschema_description:"Unique identifier for the Reddit post"`
	Summary   string `json:"summary" jsonschema_description:"Brief summary of the Reddit post content"`
	Reasoning string `json:"reasoning" jsonschema_description:"Reasoning for why this post is relevant to the search query"`
}
