package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	ctx := context.Background()
	redditClient := reddit.DefaultClient()
	posts, _, err := redditClient.Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 10,
		},
		Time: "day",
	})
	if err != nil {
		log.Fatal(err)
	}

	// for _, post := range posts {
	// 	log.WithField("title", post.Title).
	// 		WithField("id", post.ID).
	// 		WithField("score", post.Score).
	// 		WithField("created", post.Created).
	// 		Info("Post details")
	// }

	post, _, err := redditClient.Post.Get(ctx, posts[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	

	log.WithField("post", post).Info("Post details fetched by ID")
}
