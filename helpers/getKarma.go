package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) GetKarma() (int32, int32, int32) {
	future := su.System.Root.RequestFuture(su.Engine, &messages.GetKarma{
		Username: su.Username,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error getting karma: %v\n", err)
		return 0, 0, 0
	}

	if response, ok := result.(*messages.GetKarmaResponse); ok {
		if response.Success {
			log.Printf("Karma for %s - Post: %d, Comment: %d, Total: %d",
				su.Username, response.PostKarma, response.CommentKarma, response.TotalKarma)
			return response.PostKarma, response.CommentKarma, response.TotalKarma
		} else {
			log.Printf("Failed to get karma: %s", response.Error)
		}
	}
	return 0, 0, 0
}
