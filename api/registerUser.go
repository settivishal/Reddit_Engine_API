package api

import (
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) registerUser(c *gin.Context) {

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.RegisterUser{
		Username: req.Username,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		log.Printf("Error registering user as %s: %v\n", req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	response, ok := result.(*messages.RegisterUserResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		log.Printf("User %s registered successfully", req.Username)
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	} else {
		log.Printf("Error registering user %s: %s", req.Username, response.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}
