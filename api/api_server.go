package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	engine *actor.PID
	system *actor.ActorSystem
}

func NewAPIServer(engine *actor.PID, actorSystem *actor.ActorSystem) *APIServer {
	return &APIServer{
		engine: engine,
		system: actorSystem,
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)

		fmt.Printf("Request Body: %s\n", string(body))

		c.Next()
	}
}

func (s *APIServer) Start() {
	r := gin.Default()
	r.Use(RequestLogger())

	r.POST("/register", s.registerUser)     // DONE
	r.POST("/post", s.createPost)           // DONE
	r.POST("/subreddit", s.createSubreddit) // DONE
	r.POST("/join", s.joinSubreddit)        // DONE
	r.POST("/leave", s.leaveSubreddit)      // DONE
	r.POST("/vote", s.vote)                 // DONE
	r.POST("/comment", s.createComment)     // DONE
	r.GET("/karma/:username", s.getKarma)   // DONE
	// r.GET("/feed/:username", s.getFeed)
	r.POST("/message", s.sendDirectMessage)           // DONE
	r.GET("/messages/:username", s.getDirectMessages) // DONE
	r.POST("/reply", s.replyDirectMessage)            // DONE

	r.Run(":8080") // RUNS ON PORT 8080
}

// func (s *APIServer) getFeed(c *gin.Context) {
// 	username := c.Param("username")
// 	req := &messages.GetFeed{Username: username}

// 	future := s.system.Root.RequestFuture(s.engine, req, 5*time.Second)
// 	result, err := future.Result()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }
