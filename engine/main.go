package main

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"reddit_engine/helpers"
)

func main() {
	remoteConfig := remote.Configure("127.0.0.1", 8080)
	actorSystem := actor.NewActorSystem()
	remote := remote.NewRemote(actorSystem, remoteConfig)
	remote.Start()

	engineProps := actor.PropsFromProducer(func() actor.Actor {
		return &helpers.EngineActor{
			Users:            make(map[string]bool),
			Subreddits:       make(map[string]*helpers.Subreddit),
			Posts:            make(map[string]*helpers.Post),
			Comments:         make(map[string]*helpers.Comment),
			LastPostID:       0,
			LastCommentID:    0,
			PostVotes:        make(map[string]map[string]int),
			CommentVotes:     make(map[string]map[string]int),
			UserPostKarma:    make(map[string]int),
			UserCommentKarma: make(map[string]int),
			DirectMessages:   make(map[string]*helpers.DirectMessage),
      UserInbox:        make(map[string][]string),
		}
	})

	enginePID, err := actorSystem.Root.SpawnNamed(engineProps, "engine")
	if err != nil {
		log.Fatalf("Failed to spawn engine actor: %v", err)
	}

	log.Println("Engine actor PID:", enginePID)
	log.Println("Reddit Engine is running on 127.0.0.1:8080")

	select {}
}
