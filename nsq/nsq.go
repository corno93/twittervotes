package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

// TODO: READ ERRORS
var (
	pub, _ = nsq.NewProducer("localhost:4150", nsq.NewConfig())
	Sub, _ = nsq.NewConsumer("votes", "counter", nsq.NewConfig())
)

// func StartNSQ() {
// 	pub, err = nsq.NewProducer("localhost:4150", nsq.NewConfig())
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

func ConnectNSQ() {

	if err := Sub.ConnectToNSQLookupd("localhost:4161"); err != nil {
		log.Fatal(err)

	}
}

// This is the function that is called eveyrtime on the subscriber when NSQ receives something
func VotesLogic() {

}
