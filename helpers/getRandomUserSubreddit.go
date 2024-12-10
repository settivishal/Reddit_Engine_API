package helpers

// import (
// 	"math/rand"
// )

// func (su *SimulatedUser) getRandomUserSubreddit(subreddits []*Subreddit) *Subreddit {
// 	userSubreddits := make([]*Subreddit, 0)
// 	for _, subreddit := range subreddits {
// 			if subreddit.Members[su.Username] {
// 					userSubreddits = append(userSubreddits, subreddit)
// 			}
// 	}
// 	if len(userSubreddits) == 0 {
// 			return nil
// 	}
// 	return userSubreddits[rand.Intn(len(userSubreddits))]
// }