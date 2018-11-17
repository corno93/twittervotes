package main

import (
	"errors"

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

	// start mongo and read users options
	mongo.Dialdb()
	options := mongo.UsersOptions("polls")
	log.Println(options)

	// read from twitter
	ReadTwitter(options)

}
