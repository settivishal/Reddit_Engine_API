package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) createSubreddit(c *gin.Context) {

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Creator     string `json:"creator" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.CreateSubreddit{
		Name:        req.Name,
		Description: req.Description,
		Creator:     req.Creator,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subreddit"})
		return
	}

	response, ok := result.(*messages.CreateSubredditResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Subreddit created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}
