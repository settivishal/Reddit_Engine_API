package helpers

import (
	"reddit_engine/messages"
)

func (e *EngineActor) handleVote(msg *messages.Vote) *messages.VoteResponse {
	if !e.Users[msg.Voter] {
		return &messages.VoteResponse{
			Success: false,
			Error:   "User not registered",
		}
	}

	var votes map[string]map[string]int
	var karmaMap map[string]int
	var author string

	if msg.IsPost {
		post, postExists := e.Posts[msg.Id]
		if !postExists {
			return &messages.VoteResponse{
				Success: false,
				Error:   "Post does not exist",
			}
		}
		votes = e.PostVotes
		karmaMap = e.UserPostKarma
		author = post.Author
	} else {
		comment, commentExists := e.Comments[msg.Id]
		if !commentExists {
			return &messages.VoteResponse{
				Success: false,
				Error:   "Comment does not exist",
			}
		}
		votes = e.CommentVotes
		karmaMap = e.UserCommentKarma
		author = comment.Author
	}

	// Initialize vote map for the item if it doesn't exist
	if votes[msg.Id] == nil {
		votes[msg.Id] = make(map[string]int)
	}

	// Calculate vote delta
	newVote := 1
	if !msg.IsUpvote {
		newVote = -1
	}

	// Get previous vote if exists
	prevVote := votes[msg.Id][msg.Voter]

	// Update karma
	if prevVote != newVote {
		// Remove previous vote effect only if the new vote is different
		karmaMap[author] -= prevVote

		// Apply new vote
		votes[msg.Id][msg.Voter] = newVote
		karmaMap[author] += newVote
	}

	// Calculate current karma for the item
	currentKarma := 0
	for _, v := range votes[msg.Id] {
		currentKarma += v
	}

	return &messages.VoteResponse{
		Success:      true,
		CurrentKarma: int32(currentKarma),
	}
}
