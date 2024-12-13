package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) leaveSubreddit(c *gin.Context) {

	var req struct {
		Username      string `json:"username" binding:"required"`
		SubRedditName string `json:"subredditname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.LeaveSubreddit{
		Username:      req.Username,
		SubredditName: req.SubRedditName,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave subreddit"})
		return
	}

	response, ok := result.(*messages.LeaveSubredditResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Left subreddit successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}
}
