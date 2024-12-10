package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) CreateSubreddit(name, description string) bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.CreateSubreddit{
		Name:        name,
		Description: description,
		Creator:     su.Username,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error creating subreddit: %v\n", err)
		return false
	}

	if response, ok := result.(*messages.CreateSubredditResponse); ok {
		if response.Success {
			log.Printf("Successfully created subreddit: %s", name)
			return true
		} else {
			log.Printf("Failed to create subreddit: %s", response.Error)
			return false
		}
	}
	return false
}
