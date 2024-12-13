package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) createComment(c *gin.Context) {
	var req struct {
		Content         string `json:"content" binding:"required"`
		Author          string `json:"author" binding:"required"`
		PostId          string `json:"postid" binding:"required"`
		ParentCommentId string `json:"parentcommentid"` // "" if not nested
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.CreateComment{
		Content:         req.Content,
		Author:          req.Author,
		PostId:          req.PostId,
		ParentCommentId: req.ParentCommentId,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	response, ok := result.(*messages.CreateCommentResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}
}
