package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) getDirectMessages(c *gin.Context) {
	username := c.Param("username")
	req := &messages.GetDirectMessages{Username: username}

	future := s.system.Root.RequestFuture(s.engine, &messages.GetDirectMessages{
		Username: req.Username,
	}, 5*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch direct messages"})
		return
	}

	response, ok := result.(*messages.GetDirectMessagesResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	// ARRAY OF DMS FROM RECIEVED FROM A PERSON
	c.JSON(http.StatusOK, gin.H{"message": "Direct messages fetched successfully", "messages": response.Messages})

}
