package mongo

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

// global variables
var (
	Db *mgo.Session // mongo client
)

func RunMongo() {
	// read from mongo
	err := Dialdb()
	// test print something
	c := Db.DB("ballots").C("polls")
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

// function returns the options in one users collection
func UsersOptions(collection string) []string {

	c := Db.DB("ballots").C(collection)
	var pollData []PollData
	err := c.Find(nil).All(&pollData)

	if err != nil {
		//log.Printf("RunQuery : ERROR : %s\n", err)
		log.Fatalln("RunQuery : ERROR : \n", err)
	}

	// collect options in users collection
	var options []string
	for _, poll := range pollData {
		options = append(options, poll.Options...)
	}
	return options
}

// function opens the monogodb connection
func Dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	Db, err = mgo.Dial("localhost")
	Db.SetMode(mgo.Monotonic, true) //In the Monotonic consistency mode reads may not be entirely up-to-date, but they will always see the history of changes moving forward, the data read will be consistent across sequential queries in the same session, and modifications made within the session will be observed in following queries (read-your-writes).

	return err
}

// function closes mongodb
func Closedb() {
	Db.Close()
	log.Println("closed database connection")
}
