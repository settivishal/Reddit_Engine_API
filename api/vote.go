package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) vote(c *gin.Context) {
	var req struct {
		Voter    string `json:"voter" binding:"required"`
		Id       string `json:"id" binding:"required"`
		IsUpvote bool   `json:"isupvote" binding:"required"`
		IsPost   bool   `json:"ispost" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.Vote{
		Voter:    req.Voter,
		Id:       req.Id,
		IsUpvote: req.IsUpvote,
		IsPost:   req.IsPost,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	response, ok := result.(*messages.VoteResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Voted successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}
