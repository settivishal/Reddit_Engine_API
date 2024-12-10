package helpers

import (
	"log"
	"time"

	"reddit_engine/messages"
)

func (su *SimulatedUser) LeaveSubreddit(name string) bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.LeaveSubreddit{
		Username:      su.Username,
		SubredditName: name,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error leaving subreddit: %v\n", err)
		return false
	}

	if response, ok := result.(*messages.LeaveSubredditResponse); ok {
		if response.Success {
			log.Printf("Successfully left subreddit: %s", name)
			return true
		} else {
			log.Printf("Failed to leave subreddit: %s", response.Error)
			return false
		}
	}
	return false
}
