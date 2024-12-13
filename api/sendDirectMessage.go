package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) sendDirectMessage(c *gin.Context) {
	var req struct {
		FromUser string `json:"fromuser" binding:"required"`
		ToUser   string `json:"touser" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.SendDirectMessage{
		FromUser: req.FromUser,
		ToUser:   req.ToUser,
		Content:  req.Content,
	}, 5*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	response, ok := result.(*messages.SendDirectMessageResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"response": "successfully sent message", "message": response.Message})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to send direct message"})
	}
}
