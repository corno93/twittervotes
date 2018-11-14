package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PollData struct {
	ID      bson.ObjectId  `bson:"_id"`
	Title   string         `bson:"title"`
	Options []string       `bson:"options"`
	Results map[string]int `bson:"results"`
	ApiKey  string         `bson:"apikey"`
}

func main() {

}

func runMongo() {
	// read from mongo
	err := dialdb()
	// test print something
	c := db.DB("ballots").C("polls")
	var pollData []PollData
	err = c.Find(nil).All(&pollData)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}
	for _, poll := range pollData {
		log.Println(poll)
	}

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
