package helpers

import (
	"github.com/asynkron/protoactor-go/actor"
)

type SimulatedUser struct {
	Username string
	Engine   *actor.PID
	System   *actor.ActorSystem
}

type Post struct {
	ID        string
	Content   string
	Author    string
	Subreddit string
	Comments  []string
	Karma     int
	CreatedAt int
}

type Comment struct {
	ID              string
	Content         string
	Author          string
	PostID          string
	ParentCommentID string
}

type Subreddit struct {
	Name        string
	Description string
	Creator     string
	Members     map[string]bool
	Posts       []string
}

type DirectMessage struct {
	ID        string
	FromUser  string
	ToUser    string
	Content   string
	Timestamp int64
	IsRead    bool
}

type SendDirectMessage struct {
	FromUser string
	ToUser   string
	Content  string
}

type SendDirectMessageResponse struct {
	Success bool
	Message *DirectMessage
}

type GetDirectMessages struct {
	Username string
}

type GetDirectMessagesResponse struct {
	Messages []*DirectMessage
}

type EngineActor struct {
	Users            map[string]bool
	Subreddits       map[string]*Subreddit
	Posts            map[string]*Post
	Comments         map[string]*Comment
	LastPostID       int
	LastCommentID    int
	PostVotes        map[string]map[string]int // post_id -> {username: vote}
	CommentVotes     map[string]map[string]int // comment_id -> {username: vote}
	UserPostKarma    map[string]int            // username -> karma
	UserCommentKarma map[string]int            // username -> karma
	PostsByTime      []string                  // Slice to maintain posts in chronological order
	DirectMessages   map[string]*DirectMessage // message_id -> message
	UserInbox        map[string][]string       // username -> []message_ids
}
