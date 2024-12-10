package helpers

import (
	"fmt"
	"reddit_engine/messages"
)

func (e *EngineActor) handleCreateComment(msg *messages.CreateComment) *messages.CreateCommentResponse {
	if !e.Users[msg.Author] {
		return &messages.CreateCommentResponse{
			Success: false,
			Error:   "User not registered",
		}
	}

	post := e.Posts[msg.PostId]
	if post == nil {
		return &messages.CreateCommentResponse{
			Success: false,
			Error:   "Post does not exist",
		}
	}

	// Validate parent comment if provided
	if msg.ParentCommentId != "" {
		_, exists := e.Comments[msg.ParentCommentId]
		if !exists {
			return &messages.CreateCommentResponse{
				Success: false,
				Error:   "Parent comment does not exist",
			}
		}
	}

	e.LastCommentID++
	commentID := fmt.Sprintf("comment_%d", e.LastCommentID)

	comment := &Comment{
		ID:              commentID,
		Content:         msg.Content,
		Author:          msg.Author,
		PostID:          msg.PostId,
		ParentCommentID: msg.ParentCommentId,
	}

	e.Comments[commentID] = comment
	post.Comments = append(post.Comments, commentID)

	// Find and track child comments
	childCommentIds := []string{}
	for _, existingCommentID := range post.Comments {
		existingComment := e.Comments[existingCommentID]
		if existingComment.ParentCommentID == msg.ParentCommentId {
			childCommentIds = append(childCommentIds, existingCommentID)
		}
	}

	return &messages.CreateCommentResponse{
		Success:         true,
		CommentId:       commentID,
		ChildCommentIds: childCommentIds,
	}
}
