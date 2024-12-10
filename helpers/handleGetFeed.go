package helpers

import (
		"reddit_engine/messages"
)

func (e *EngineActor) handleGetFeed(msg *messages.GetFeed) *messages.GetFeedResponse {
	var feedPosts []*messages.Post
	var postsToConsider []string

	if msg.Subreddit != "" {
			// Get subreddit-specific feed
			subreddit, exists := e.Subreddits[msg.Subreddit]
			if !exists {
					return &messages.GetFeedResponse{
							Success: false,
							Error: "Subreddit does not exist",
					}
			}
			postsToConsider = subreddit.Posts
	} else {
			// Get global feed
			postsToConsider = e.PostsByTime
	}

	// Calculate pagination bounds
	start := int(msg.Offset)
	end := start + int(msg.Limit)
	if end > len(postsToConsider) {
			end = len(postsToConsider)
	}
	if start >= end {
			return &messages.GetFeedResponse{
					Success: true,
					Posts: []*messages.Post{},
			}
	}

	// Get posts within pagination bounds
	for _, postID := range postsToConsider[start:end] {
			post := e.Posts[postID]
			karma := 0
			if votes := e.PostVotes[postID]; votes != nil {
					for _, vote := range votes {
							karma += vote
					}
			}

			feedPosts = append(feedPosts, &messages.Post{
					Id: post.ID,
					Content: post.Content,
					Author: post.Author,
					Subreddit: post.Subreddit,
					Karma: int32(karma),
					Comments: post.Comments,
			})
	}

	return &messages.GetFeedResponse{
			Success: true,
			Posts: feedPosts,
	}
}