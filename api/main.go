package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/corno93/twittervotes/mongo"
)

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
)


type contextKey struct {
	name string
}

//we are going to add a key to store our API key value in
var contextKeyAPIKey = &contextKey{"api-key"}

//  helper that will, given a context, extract the key
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

// The same-origin security policy mandates that AJAX requests in web browsers be allowed only for services hosted on the same domain
// The CORS (Cross-origin resource sharing) technique circumnavigates the same-origin policy, allowing us to build a service capable of serving websites hosted on other domains
// it just sets the appropriate header on the ResponseWriter type and calls the specified http.HandlerFunc type.
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers",
			"Location")
		fn(w, r)
	}
}

// Server is the API server.  encapsulates all the dependencies for our handlers and construct it with a database connection
// Our handler functions will be methods of this server, which is how they will be able to access the database session.
type Server struct {
	db *mgo.Session
}
