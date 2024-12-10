package helpers

import (
	"reddit_engine/messages"
)

func (e *EngineActor) handleGetKarma(msg *messages.GetKarma) *messages.GetKarmaResponse {
	if !e.Users[msg.Username] {
		return &messages.GetKarmaResponse{
			Success: false,
			Error:   "User not registered",
		}
	}

	postKarma := e.UserPostKarma[msg.Username]
	commentKarma := e.UserCommentKarma[msg.Username]

	return &messages.GetKarmaResponse{
		Success:      true,
		PostKarma:    int32(postKarma),
		CommentKarma: int32(commentKarma),
		TotalKarma:   int32(postKarma + commentKarma),
	}
}
