package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/corno93/twittervotes/mongo"

	"log"

	"gopkg.in/mgo.v2/bson"
)

/*
	Project: Twitter tracking project

	File: twitter votes

	Dependencies: twitter secrets must be saved as environment variables

	Description:
		- collect user's options from mongo
		- track user's options on twitter
		- publish tracked words via NSQ bus

*/

// Errors
var (
	errFailedAuth = errors.New("Could not read all authorisation environment variables")
)

type PollData struct {
	ID      bson.ObjectId  `bson:"_id"`
	Title   string         `bson:"title"`
	Options []string       `bson:"options"`
	Results map[string]int `bson:"results"`
	ApiKey  string         `bson:"apikey"`
}

func main() {

	// create channels
	votes := make(chan string) // chan for votes
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start mongo and read users options
	mongo.Dialdb()
	options := mongo.UsersOptions("polls")
	log.Println(options)

	// read from twitter
	go ReadTwitter(votes, options)

	// send on NSQ
	//	go nsq.PublishVotes(votes)

	for {
		select {
		case vote := <-votes:
			log.Println("received ", vote)
		case shutdown := <-signalChan:
			log.Println("System Shutdown", shutdown)
			//	CloseTwitterConn()
			ShutDownTwitter()
			break
			//	nsq.ShutDownNSQ()
		}
	}

	log.Println("Dead")

}
