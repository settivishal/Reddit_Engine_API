package helpers

import (
	"log"
	"time"

	"reddit_engine/messages"
)

func (su *SimulatedUser) Vote(itemID string, isUpvote bool, isPost bool) bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.Vote{
		Voter:    su.Username,
		Id:       itemID,
		IsUpvote: isUpvote,
		IsPost:   isPost,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error voting: %v\n", err)
		return false
	}

	if response, ok := result.(*messages.VoteResponse); ok {
		if response.Success {
			voteType := "upvoted"
			if !isUpvote {
				voteType = "downvoted"
			}
			itemType := "post"
			if !isPost {
				itemType = "comment"
			}
			log.Printf("Successfully %s %s %s. Current karma: %d",
				voteType, itemType, itemID, response.CurrentKarma)
			return true
		} else {
			// log.Printf("Failed to vote: %s", response.Error)
			return false
		}
	}
	return false
}
