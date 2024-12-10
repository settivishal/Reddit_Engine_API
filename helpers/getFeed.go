package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) GetFeed(subreddit string, limit int32, offset int32) []*messages.Post {
	future := su.System.Root.RequestFuture(su.Engine, &messages.GetFeed{
			Subreddit: subreddit,
			Limit: limit,
			Offset: offset,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
			log.Printf("Error getting feed: %v\n", err)
			return nil
	}

	if response, ok := result.(*messages.GetFeedResponse); ok {
			if response.Success {
					log.Printf("Retrieved %d posts from feed", len(response.Posts))
					return response.Posts
			} else {
					log.Printf("Failed to get feed: %s", response.Error)
			}
	}
	return nil
}