package helpers

import (
	"log"
	"reddit_engine/messages"
	"time"
)

func (su *SimulatedUser) Register() bool {
	future := su.System.Root.RequestFuture(su.Engine, &messages.RegisterUser{
		Username: su.Username,
	}, 10*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error registering user %s: %v\n", su.Username, err)
		return false
	}

	response, ok := result.(*messages.RegisterUserResponse)
	if !ok {
		log.Printf("Invalid response type: %v", result)
		return false
	}

	if response.Success {
		log.Printf("User %s registered successfully", su.Username)
		return true
	} else {
		log.Printf("Error registering user %s: %s", su.Username, response.Error)
		return false
	}
}
