package main

import (
	"context"
	"flag"
	"log"
	"net/http"

<<<<<<< HEAD
	"github.com/corno93/twittervotes/mongo"
	mgo "gopkg.in/mgo.v2"
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
=======
	mgo "gopkg.in/mgo.v2"
)


func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	defer db.Close()
	s := &Server{
		db: db}
	// which is a request multiplexer provided by the Go standard library, and register a single handler for all requests that begin with
	// the /polls/ path. Note that the handlePolls handler is a method on our server, and this is how it will be able to access the database.
>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
	log.Println("Starting web server on", *addr)
	http.ListenAndServe(":8080", mux)
	log.Println("Stopping...")
}

<<<<<<< HEAD
=======
//Note that contextKey and contextKeyAPIKey are internal (they start with a lowercase letter) but APIKey will be exported
>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
type contextKey struct {
	name string
}

<<<<<<< HEAD
// add a key to store our API key value in
=======
//we are going to add a key to store our API key value in
>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
var contextKeyAPIKey = &contextKey{"api-key"}

//  helper that will, given a context, extract the key
func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

<<<<<<< HEAD
// withAPIKey recieves and returns a http.HandlerFunc type. It validates the API key
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {

=======
// withAPIKey function both takes an http.HandlerFunc type as an argument and returns one; this is what we mean by wrapping in this context.
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {

	// return function performs a check for the key query parameter by calling isValidAPIKey
>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
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
<<<<<<< HEAD
=======
// we are handling CORS explicitly so we can understand exactly what is going on; for real production code, you should consider employing an open source solution, such as https://github.com/faster ness/cors.

>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
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
