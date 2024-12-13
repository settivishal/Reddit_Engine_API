package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) createPost(c *gin.Context) {

	var req struct {
		Content   string `json:"content" binding:"required"`
		Author    string `json:"author" binding:"required"`
		Subreddit string `json:"subreddit" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.CreatePost{
		Content:   req.Content,
		Author:    req.Author,
		Subreddit: req.Subreddit,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	response, ok := result.(*messages.CreatePostResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}
