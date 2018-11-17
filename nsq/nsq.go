package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

var (
	pub, _ = nsq.NewProducer("localhost:4150", nsq.NewConfig())
)

// func StartNSQ() {
// 	pub, _ = nsq.NewProducer("localhost:4150", nsq.NewConfig())
// }

func PublishVotes(votes chan string) {
	for vote := range votes {
		pub.Publish("votes", []byte(vote)) // publish vote
	}
}

func ShutDownNSQ() {
	log.Println("Publisher: Stopping")
	pub.Stop()
	log.Println("Publisher: Stopped")
}
