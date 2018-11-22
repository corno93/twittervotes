package tests

import (
	"log"
	"testing"

	"github.com/corno93/twittervotes/mongo"
)

//This test ensures we can connect to the local/remote database and that we can retrive data from it
func TestDB(t *testing.T) {

	err := mongo.Dialdb()
	if err != nil {
		t.Error("There should be no error")
		log.Println(err)
	}
	options := mongo.UsersOptions("polls")
	log.Println(options)
	if options == nil {
		t.Error("No data in 'polls' collection")
	}

}
