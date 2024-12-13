package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) replyDirectMessage(c *gin.Context) {
	var req struct {
		FromUser         string `json:"fromuser" binding:"required"`
		ReplyToMessageId string `json:"replytomessageid" binding:"required"`
		Content          string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.ReplyDirectMessage{
		FromUser:         req.FromUser,
		ReplyToMessageId: req.ReplyToMessageId,
		Content:          req.Content,
	}, 5*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	response, ok := result.(*messages.ReplyDirectMessageResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"response": "successfully replied to message", "message": response.Message})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to reply to direct message"})
	}

}
