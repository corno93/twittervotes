package tests

import (
	"testing"

	"github.com/corno93/twittervotes/twitter"
)

func TestTwitter(t *testing.T) {

	// read and test env variables (check cred and authClient)
	twitter.SetupTwitterAuth()

}
