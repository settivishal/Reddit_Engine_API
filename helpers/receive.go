package helpers

import (
	"fmt"
	"log"
	"reddit_engine/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func (e *EngineActor) Receive(context actor.Context) {
	log.Printf("Engine received message: %T", context.Message())

	switch msg := context.Message().(type) {

	case *messages.Handshake:
		context.Respond(&messages.HandshakeResponse{
			Success: true,
			Message: "Engine is up and ready to communicate!",
		})

	case *messages.RegisterUser:
		if _, exists := e.Users[msg.Username]; exists {
			context.Respond(&messages.RegisterUserResponse{
				Success: false,
				Error:   "User already exists",
			})
		} else {
			e.Users[msg.Username] = true
			context.Respond(&messages.RegisterUserResponse{Success: true})
			log.Printf("User registered: %s\n", msg.Username)
		}

	case *messages.CreatePost:
		if !e.Users[msg.Author] {
			context.Respond(&messages.CreatePostResponse{
				Success: false,
				Error:   "User not registered",
			})
			return
		}

		subreddit, exists := e.Subreddits[msg.Subreddit]
		if !exists {
			context.Respond(&messages.CreatePostResponse{
				Success: false,
				Error:   "Subreddit does not exist",
			})
			return
		}

		// Check if user is a member of the subreddit
		if !subreddit.Members[msg.Author] {
			context.Respond(&messages.CreatePostResponse{
				Success: false,
				Error:   "User is not a member of this subreddit",
			})
			return
		}

		e.LastPostID++
		postID := fmt.Sprintf("post_%d", e.LastPostID)

		// Create a new Post struct instance
		post := &Post{
			ID:        postID,
			Content:   msg.Content,
			Author:    msg.Author,
			Subreddit: msg.Subreddit,
			Comments:  make([]string, 0),
		}

		e.Posts[postID] = post
		subreddit.Posts = append(subreddit.Posts, postID)

		context.Respond(&messages.CreatePostResponse{
			Success: true,
			PostId:  postID,
		})
		log.Printf("Post created in %s by %s with ID: %s\n", msg.Subreddit, msg.Author, postID)

		e.PostsByTime = append(e.PostsByTime, postID)

	case *messages.CreateSubreddit:
		if !e.Users[msg.Creator] {
			context.Respond(&messages.CreateSubredditResponse{
				Success: false,
				Error:   "User not registered",
			})
			return
		}

		if _, exists := e.Subreddits[msg.Name]; exists {
			context.Respond(&messages.CreateSubredditResponse{
				Success: false,
				Error:   "Subreddit already exists",
			})
			return
		}

		e.Subreddits[msg.Name] = &Subreddit{
			Name:        msg.Name,
			Description: msg.Description,
			Creator:     msg.Creator,
			Members:     map[string]bool{msg.Creator: true}, // Creator automatically joins
			Posts:       make([]string, 0),
		}

		context.Respond(&messages.CreateSubredditResponse{Success: true})
		log.Printf("Subreddit created: %s by %s\n", msg.Name, msg.Creator)

	case *messages.JoinSubreddit:
		if !e.Users[msg.Username] {
			context.Respond(&messages.JoinSubredditResponse{
				Success: false,
				Error:   "User not registered",
			})
			return
		}

		subreddit, exists := e.Subreddits[msg.SubredditName]
		if !exists {
			context.Respond(&messages.JoinSubredditResponse{
				Success: false,
				Error:   "Subreddit does not exist",
			})
			return
		}

		if subreddit.Members[msg.Username] {
			context.Respond(&messages.JoinSubredditResponse{
				Success: false,
				Error:   "User is already a member",
			})
			return
		}

		subreddit.Members[msg.Username] = true
		context.Respond(&messages.JoinSubredditResponse{Success: true})
		log.Printf("User %s joined subreddit: %s\n", msg.Username, msg.SubredditName)

	case *messages.LeaveSubreddit:
		if !e.Users[msg.Username] {
			context.Respond(&messages.LeaveSubredditResponse{
				Success: false,
				Error:   "User not registered",
			})
			return
		}

		subreddit, exists := e.Subreddits[msg.SubredditName]
		if !exists {
			context.Respond(&messages.LeaveSubredditResponse{
				Success: false,
				Error:   "Subreddit does not exist",
			})
			return
		}

		if !subreddit.Members[msg.Username] {
			context.Respond(&messages.LeaveSubredditResponse{
				Success: false,
				Error:   "User is not a member",
			})
			return
		}

		if msg.Username == subreddit.Creator {
			context.Respond(&messages.LeaveSubredditResponse{
				Success: false,
				Error:   "Creator cannot leave their subreddit",
			})
			return
		}

		delete(subreddit.Members, msg.Username)
		context.Respond(&messages.LeaveSubredditResponse{Success: true})
		log.Printf("User %s left subreddit: %s\n", msg.Username, msg.SubredditName)

	case *messages.Vote:
		context.Respond(e.handleVote(msg))

	case *messages.CreateComment:
		context.Respond(e.handleCreateComment(msg))

	case *messages.GetKarma:
		context.Respond(e.handleGetKarma(msg))

	case *messages.GetFeed:
		context.Respond(e.handleGetFeed(msg))

	case *messages.SendDirectMessage:
		localMsg := &SendDirectMessage{
			FromUser: msg.FromUser,
			ToUser:   msg.ToUser,
			Content:  msg.Content,
		}
		response, err := e.HandleSendDirectMessage(localMsg)
		if err != nil {
			context.Respond(&messages.SendDirectMessageResponse{Success: false})
			return
		}
		// Convert local response to protobuf response
		protoResponse := &messages.SendDirectMessageResponse{
			Success: response.Success,
			Message: &messages.DirectMessage{
				Id:        response.Message.ID,
				FromUser:  response.Message.FromUser,
				ToUser:    response.Message.ToUser,
				Content:   response.Message.Content,
				Timestamp: response.Message.Timestamp,
				IsRead:    response.Message.IsRead,
			},
		}
		context.Respond(protoResponse)

	case *messages.GetDirectMessages:
		localMsg := &GetDirectMessages{
			Username: msg.Username,
		}
		response, err := e.HandleGetDirectMessages(localMsg)
		if err != nil {
			context.Respond(&messages.GetDirectMessagesResponse{})
			return
		}
		// Convert local response to protobuf response
		protoMessages := make([]*messages.DirectMessage, len(response.Messages))
		for i, dm := range response.Messages {
			protoMessages[i] = &messages.DirectMessage{
				Id:        dm.ID,
				FromUser:  dm.FromUser,
				ToUser:    dm.ToUser,
				Content:   dm.Content,
				Timestamp: dm.Timestamp,
				IsRead:    dm.IsRead,
			}
		}
		context.Respond(&messages.GetDirectMessagesResponse{
			Messages: protoMessages,
		})

	case *messages.ReplyDirectMessage:
		response, err := e.HandleReplyDirectMessage(msg)
		if err != nil {
			context.Respond(&messages.ReplyDirectMessageResponse{Success: false})
			return
		}
		context.Respond(response)
	default:
		log.Printf("Unknown message received: %T\n", msg)
	}
}
