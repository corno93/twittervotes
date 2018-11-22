package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/corno93/twittervotes/mongo"
	"github.com/corno93/twittervotes/nsq"

	"log"
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

func main() {

	// create channels
	votes := make(chan string)            // chan for votes
	signalChan := make(chan os.Signal, 1) // gracefully shutdown when cntl-c is hit
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start mongo and read users options
	mongo.Dialdb()
	options := mongo.UsersOptions("polls")

	// read from twitter
	go ReadTwitter(votes, options)

	// send on NSQ
	go nsq.PublishVotes(votes)

readChannels:
	for {
		select {
		case vote := <-votes:
			log.Println("received ", vote)
		case shutdown := <-signalChan:
			log.Println("System Shutdown", shutdown)
			ShutDownTwitter()
			mongo.Closedb()
			close(votes)
			nsq.ShutDownNSQ()
			break readChannels

		}
	}
	log.Println("System dead")

}
