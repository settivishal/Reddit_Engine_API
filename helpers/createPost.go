package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) CreatePost(content, subreddit string) bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.CreatePost{
		Content:   content,
		Author:    su.Username,
		Subreddit: subreddit,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error creating post: %v\n", err)
		return false
	}

	if response, ok := result.(*messages.CreatePostResponse); ok {
		if response.Success {
			log.Printf("Post created successfully in %s with ID: %s", subreddit, response.PostId)
			return true
		} else {
			log.Printf("Failed to create post: %s", response.Error)
			return false
		}
	}
	return false
}
