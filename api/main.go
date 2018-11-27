package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/corno93/twittervotes/mongo"
	mgo "gopkg.in/mgo.v2"
)

// Server struct will encapsulates all the dependencies for our handlers and include a database connection
type Server struct {
	db *mgo.Session
}

func main() {
	var (
		addr = flag.String("addr", ":8080", "endpoint address")
	)

	// dial into mongo
	mongo.Dialdb()

	s := &Server{
		db: mongo.Db}

	// configure multiplexer with handlePolls handler using CORS technique and an API key for good practice
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
	log.Println("Starting web server on", *addr)
	http.ListenAndServe(":8080", mux)
	log.Println("Stopping...")
}

// add a key to store our API key value in
type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

// helper that will, given a context, extract the key
func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

// withAPIKey recieves and returns a http.HandlerFunc type. It validates the API key
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {

	// return function performs a check for the key query parameter by calling isValidAPIKey
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		// otherwise we put the key into the context and call the next handler
		ctx := context.WithValue(r.Context(),
			contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

// helper function that users CORS headers
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // share response
		w.Header().Set("Access-Control-Expose-Headers",    // expose location
			"Location")
		fn(w, r)
	}
}
