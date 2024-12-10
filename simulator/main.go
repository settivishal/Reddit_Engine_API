package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"sync"

	"reddit_engine/helpers"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

const (
	numUsers                 = 10000
	numSubreddits            = 500
	avgSubredditMembers      = 20
	zipfAlpha                = 1.5
	avgPostsPerUser          = 20
	avgRepostsPerUser        = 10
	avgCommentsPerUser       = 20
	avgVotesPerUser          = 100
	connectionPeriod         = 5 * time.Second
	avgDirectMessagesPerUser = 10
	avgDirectMessageReplies  = 10
)

func getRandomUserSubreddit(user *helpers.SimulatedUser, subreddits []*helpers.Subreddit) *helpers.Subreddit {
	userSubreddits := make([]*helpers.Subreddit, 0)
	for _, subreddit := range subreddits {
		if subreddit.Members[user.Username] {
			userSubreddits = append(userSubreddits, subreddit)
		}
	}
	if len(userSubreddits) == 0 {
		return nil
	}
	return userSubreddits[rand.Intn(len(userSubreddits))]
}

// Helper function to get random valid post from a subreddit
func getRandomPost(subreddit *helpers.Subreddit) string {
	if len(subreddit.Posts) == 0 {
		return ""
	}
	return subreddit.Posts[rand.Intn(len(subreddit.Posts))]
}

func simulateUsers(system *actor.ActorSystem, enginePID *actor.PID) {
	// Create users and subreddits
	users := make([]*helpers.SimulatedUser, numUsers)
	subreddits := make([]*helpers.Subreddit, numSubreddits)

	// register users to engine actor
	for i := range users {
		users[i] = &helpers.SimulatedUser{
			Username: fmt.Sprintf("user%d", i),
			Engine:   enginePID,
			System:   system,
		}
		users[i].Register()
	}

	// create subreddits, random user will create a subreddit
	for i := range subreddits {
		creator := users[rand.Intn(numUsers)]

		subreddits[i] = &helpers.Subreddit{
			Name:        fmt.Sprintf("subreddit%d", i),
			Description: fmt.Sprintf("Description for subreddit %d", i),
			Creator:     creator.Username,
			Members:     make(map[string]bool),
			Posts:       make([]string, 0),
		}

		creator.CreateSubreddit(subreddits[i].Name, subreddits[i].Description)
		// Add creator as first member
		subreddits[i].Members[creator.Username] = true
		creator.JoinSubreddit(subreddits[i].Name)
	}

	zipf := helpers.ZipfDistribution(numUsers, zipfAlpha)

	// Assign subreddit members using Zipf distribution
	for i, user := range users {
		numSubredditsToJoin := zipf[i]

		// Create a list of subreddits the user isn't already a member of
		availableSubreddits := make([]int, 0)
		for j := range subreddits {
			if !subreddits[j].Members[user.Username] {
				availableSubreddits = append(availableSubreddits, j)
			}
		}

		// Join random subreddits
		for j := 0; j < numSubredditsToJoin && len(availableSubreddits) > 0; j++ {
			// Pick a random index from available subreddits
			idx := rand.Intn(len(availableSubreddits))
			subredditID := availableSubreddits[idx]

			// Join the subreddit
			user.JoinSubreddit(subreddits[subredditID].Name)
			subreddits[subredditID].Members[user.Username] = true

			// Remove this subreddit from available list
			availableSubreddits = append(availableSubreddits[:idx], availableSubreddits[idx+1:]...)
		}
	}

	fmt.Println("Subreddits created and members assigned")

	// Simulate user connections and disconnections
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user *helpers.SimulatedUser) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				// Simulate connection period
				time.Sleep(time.Second)

				// Simulate post creation
				numPosts := rand.Intn(avgPostsPerUser)
				for j := 0; j < numPosts; j++ {
					subreddit := getRandomUserSubreddit(user, subreddits)
					if subreddit == nil {
						continue
					}
					postID := fmt.Sprintf("post%d_%d_%s", i, j, user.Username)
					success := user.CreatePost(postID, subreddit.Name)
					if success {
						subreddit.Posts = append(subreddit.Posts, postID)
						fmt.Printf("User %s created post %s in subreddit %s\n",
							user.Username, postID, subreddit.Name)
					} else {
						fmt.Printf("Failed to create post %s in subreddit %s\n",
							postID, subreddit.Name)
					}
				}

				time.Sleep(connectionPeriod)

				// Simulate Commenting
				numComments := rand.Intn(avgCommentsPerUser)
				for j := 0; j < numComments; j++ {
					subreddit := getRandomUserSubreddit(user, subreddits)
					if subreddit == nil {
						continue
					}
					validPosts := subreddit.Posts
					if len(validPosts) == 0 {
						continue
					}

					// Select a random post from valid posts
					randPost := rand.Intn(len(validPosts))
					postID := validPosts[randPost]
					commentContent := fmt.Sprintf("Comment #%d from %s on post %s", j, user.Username, postID)

					_, err := user.CreateComment(commentContent, postID, "")
					if err == nil {
						fmt.Printf("User %s commented on post %s in subreddit %s\n",
							user.Username, postID, subreddit.Name)
					} else {
						fmt.Printf("Failed to create comment: %v\n", err)
					}

					time.Sleep(50 * time.Millisecond)
				}

				// Simulate Voting
				numVotes := rand.Intn(avgVotesPerUser)
				// votedPosts := make(map[string]bool) // Track which posts user has already voted on
				for j := 0; j < numVotes; j++ {
					subreddit := getRandomUserSubreddit(user, subreddits)
					if subreddit == nil {
						continue
					}

					validPosts := subreddit.Posts
					if len(validPosts) == 0 {
						continue
					}

					randPost := rand.Intn(len(validPosts))
					postID := validPosts[randPost]

					isUpvote := rand.Intn(2) == 0
					err := user.Vote(postID, isUpvote, true)
					if !err {
						voteType := "upvoted"
						if !isUpvote {
							voteType = "downvoted"
						}
						fmt.Printf("User %s %s post %s in subreddit %s\n",
							user.Username, voteType, postID, subreddit.Name)
					}

					time.Sleep(50 * time.Millisecond)
				}

				// Simulate direct messaging
				numDirectMessages := rand.Intn(avgDirectMessagesPerUser)
				for j := 0; j < numDirectMessages; j++ {
					// Pick a random recipient that isn't the sender
					recipientIndex := rand.Intn(numUsers)
					if users[recipientIndex].Username == user.Username {
						continue
					}

					// Send direct message
					content := fmt.Sprintf("Hello from %s! Message #%d", user.Username, j)
					success := user.SendDirectMessage(users[recipientIndex].Username, content)

					if success != nil && success.Success {
						fmt.Printf("User %s sent direct message to %s\n",
							user.Username, users[recipientIndex].Username)

						// Simulate replies
						numReplies := rand.Intn(avgDirectMessageReplies)
						for k := 0; k < numReplies; k++ {
							replyContent := fmt.Sprintf("Reply #%d to message from %s", k, user.Username)
							replySuccess := users[recipientIndex].ReplyDirectMessage(content, replyContent)

							if replySuccess {
								fmt.Printf("User %s replied to message from %s\n",
									users[recipientIndex].Username, user.Username)
							}

							// Simulate reading messages
							user.GetDirectMessages()
							users[recipientIndex].GetDirectMessages()

							time.Sleep(50 * time.Millisecond)
						}
					} else {
						fmt.Printf("Failed to send direct message from %s to %s\n",
							user.Username, users[recipientIndex].Username)
					}
					time.Sleep(100 * time.Millisecond)
				}

				// Simulate getting karma
				user.GetKarma()

				// Simulate getting feed
				user.GetFeed("", 10, 0)

				// Simulate disconnection period
				time.Sleep(time.Second)
				user.LeaveSubreddit(subreddits[rand.Intn(numSubreddits)].Name)
				// fmt.Printf("User %s disconnected\n", user.Username)
			}
		}(user)
	}

	wg.Wait()
	log.Println("Simulation finished")
}

func main() {
	actorSystem := actor.NewActorSystem()
	remoteConfig := remote.Configure("127.0.0.1", 8081)
	remoteContext := remote.NewRemote(actorSystem, remoteConfig)
	remoteContext.Start()

	log.Println("Simulator remote server started on 127.0.0.1:8081")

	enginePID := actor.NewPID("127.0.0.1:8080", "engine")
	log.Printf("Engine PID: %v", enginePID)

	// Run simulations
	simulateUsers(actorSystem, enginePID)
}
