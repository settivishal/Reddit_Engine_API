package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) getKarma(c *gin.Context) {
	username := c.Param("username")
	req := &messages.GetKarma{Username: username}

	future := s.system.Root.RequestFuture(s.engine, &messages.GetKarma{
		Username: req.Username,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch karma"})
		return
	}

	response, ok := result.(*messages.GetKarmaResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Karma fetched successfully", "total": response.TotalKarma})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}
