package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/corno93/twittervotes/mongo"
	"github.com/nsqio/go-nsq"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	counts     map[string]int
	countsLock sync.Mutex
)

const updateDuration = 1 * time.Second

func main() {

	// start mongo
	mongo.Dialdb()

	defer func() {
		log.Println("Closing database connection...")
		mongo.Closedb()
	}()

	// get pointer to the users poll data
	pollData := mongo.Db.DB("ballots").C("polls")

	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
		return
	}

	// handler function that is called everytime nsq receives something
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		log.Fatal(err)
		return
	}

	// configure ticker interrupt so update database.
	// use cntr c to gracefully stop script
	log.Println("Waiting for votes on nsq...")
	ticker := time.NewTicker(updateDuration)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-ticker.C:
			doCount(&countsLock, &counts, pollData)
		case <-termChan:
			ticker.Stop()
			q.Stop()
		case <-q.StopChan:
			// finished
			return
		}
	}

}

func doCount(countsLock *sync.Mutex, counts *map[string]int, pollData *mgo.Collection) {
	countsLock.Lock()
	defer countsLock.Unlock()
	if len(*counts) == 0 {
		log.Println("No new votes, skipping database update")
		return
	}
	log.Println("Updating database...")
	log.Println(*counts)
	ok := true
	for option, count := range *counts {
		sel := bson.M{"options": bson.M{"$in": []string{option}}}
		up := bson.M{"$inc": bson.M{"results." + option: count}}
		if _, err := pollData.UpdateAll(sel, up); err != nil {
			log.Println("failed to update:", err)
			ok = false
		}
	}
	if ok {
		log.Println("Finished updating database...")
		*counts = nil // reset counts
	}
}
