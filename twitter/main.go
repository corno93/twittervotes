package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
	Project: Distributed project example

	File: twitter votes

	Dependencies: this package requies the setup.sh file to be in the same directory path

	Description: this package is responsible for collecting reading user's
		options in the mongodb, collecting these words via the twitter api
		and finally pushing these collected words on nsq to go to our count package.

*/

// global variabels
var (
	db *mgo.Session // mongo client
)

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

	// read authorisation keys
	// err := readAuth()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(auth)

	// read from mongo
	err := dialdb()
	if err != nil {
		log.Fatalln(err)
	}
	results := returnResults("polls")
	fmt.Println(results)

	// read from twitter
	readTwitter()

}

// read from twitter
func readTwitter() {

	// make the url
	u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")

	query := make(url.Values)
	options := []string{"c", "go"}
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("creating filter request failed:", err)
	}

}

// return results for a collection
func returnResults(collection string) []PollData {

	c := db.DB("ballots").C(collection)
	var pollData []PollData
	err := c.Find(nil).All(&pollData)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
	return pollData

}

// function opens the monogodb connection
func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost")
	db.SetMode(mgo.Monotonic, true) //In the Monotonic consistency mode reads may not be entirely up-to-date, but they will always see the history of changes moving forward, the data read will be consistent across sequential queries in the same session, and modifications made within the session will be observed in following queries (read-your-writes).

	return err
}

// function closes mongodb
func closedb() {
	db.Close()
	log.Println("closed database connection")
}
