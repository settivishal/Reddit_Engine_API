package helpers

import (
	"time"
	"fmt"
"math/rand"
	"reddit_engine/messages"
)
func (u *SimulatedUser) SendDirectMessage(toUser string, content string) *SendDirectMessageResponse {
	response, err := u.System.Root.RequestFuture(u.Engine, &messages.SendDirectMessage{
			FromUser: u.Username,
			ToUser:   toUser,
			Content:  content,
	}, 5*time.Second).Result()
	
	if err != nil {
			fmt.Printf("Error sending direct message from %s to %s: %v\n", u.Username, toUser, err)
			return nil
	}
	
	if resp, ok := response.(*SendDirectMessageResponse); ok {
			fmt.Printf("User %s sent direct message to %s\n", u.Username, toUser)
			return resp
	}
	return nil
}

func (u *SimulatedUser) ReplyDirectMessage(messageID string, content string) bool {
	_, err := u.System.Root.RequestFuture(u.Engine, &messages.ReplyDirectMessage{
			FromUser:         u.Username,
			ReplyToMessageId: messageID,
			Content:         content,
	}, 5*time.Second).Result()
	
	if err != nil {
			fmt.Printf("Error replying to message %s from %s: %v\n", messageID, u.Username, err)
			return false
	}
	
	fmt.Printf("User %s replied to message %s\n", u.Username, messageID)
	return true
}

func (u *SimulatedUser) GetDirectMessages() {
	_, err := u.System.Root.RequestFuture(u.Engine, &messages.GetDirectMessages{
			Username: u.Username,
	}, 5*time.Second).Result()
	
	if err != nil {
			fmt.Printf("Error getting messages for user %s: %v\n", u.Username, err)
			return
	}
	
	fmt.Printf("User %s checked their messages\n", u.Username)
}



func (su *SimulatedUser) getRandomPost(subreddit *Subreddit) string {
	if len(subreddit.Posts) == 0 {
			return ""
	}
	return subreddit.Posts[rand.Intn(len(subreddit.Posts))]
}

func (su *SimulatedUser) getRandomUserSubreddit(subreddits []*Subreddit) *Subreddit {
	userSubreddits := make([]*Subreddit, 0)
	for _, subreddit := range subreddits {
			if subreddit.Members[su.Username] {
					userSubreddits = append(userSubreddits, subreddit)
			}
	}
	if len(userSubreddits) == 0 {
			return nil
	}
	return userSubreddits[rand.Intn(len(userSubreddits))]
}