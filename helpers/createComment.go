package helpers

import (
	"log"
	"time"

	"reddit_engine/messages"
)

func (su *SimulatedUser) CreateComment(content, postID, parentCommentID string) (string, []string) {
	future := su.System.Root.RequestFuture(su.Engine, &messages.CreateComment{
		Content:         content,
		Author:          su.Username,
		PostId:          postID,
		ParentCommentId: parentCommentID,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error creating comment: %v\n", err)
		return "", nil
	}

	if response, ok := result.(*messages.CreateCommentResponse); ok {
		if response.Success {
			log.Printf("Comment created successfully with ID: %s", response.CommentId)
			log.Printf("Child comment IDs: %v", response.ChildCommentIds)
			return response.CommentId, response.ChildCommentIds
		} else {
			log.Printf("Failed to create comment: %s", response.Error)
			return "", nil
		}
	}
	return "", nil
}
