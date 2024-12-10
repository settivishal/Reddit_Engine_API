package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) JoinSubreddit(name string) bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.JoinSubreddit{
		Username:      su.Username,
		SubredditName: name,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error joining subreddit: %v\n", err)
		return false
	}

	if response, ok := result.(*messages.JoinSubredditResponse); ok {
		if response.Success {
			log.Printf("Successfully joined subreddit: %s", name)
			return true
		} else {
			log.Printf("Failed to join subreddit: %s", response.Error)
			return false
		}
	}
	return false
}
