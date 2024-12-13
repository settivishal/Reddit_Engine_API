package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reddit_engine/messages"
	"time"

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

func (s *APIServer) createPost(c *gin.Context) {

	var req struct {
		Content   string `json:"content" binding:"required"`
		Author    string `json:"author" binding:"required"`
		Subreddit string `json:"subreddit" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.CreatePost{
		Content:   req.Content,
		Author:    req.Author,
		Subreddit: req.Subreddit,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	response, ok := result.(*messages.CreatePostResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}

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

func (s *APIServer) joinSubreddit(c *gin.Context) {

	var req struct {
		Username      string `json:"username" binding:"required"`
		SubRedditName string `json:"subredditname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.JoinSubreddit{
		Username:      req.Username,
		SubredditName: req.SubRedditName,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join subreddit"})
		return
	}

	response, ok := result.(*messages.JoinSubredditResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Joined subreddit successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}

}

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

func (s *APIServer) createComment(c *gin.Context) {
	var req struct {
		Content         string `json:"content" binding:"required"`
		Author          string `json:"author" binding:"required"`
		PostId          string `json:"postid" binding:"required"`
		ParentCommentId string `json:"parentcommentid"` // "" if not nested
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	future := s.system.Root.RequestFuture(s.engine, &messages.CreateComment{
		Content:         req.Content,
		Author:          req.Author,
		PostId:          req.PostId,
		ParentCommentId: req.ParentCommentId,
	}, 10*time.Second)

	result, err := future.Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	response, ok := result.(*messages.CreateCommentResponse)

	if !ok {
		log.Printf("Invalid response type: %v", result)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from server"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
	}
}

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
