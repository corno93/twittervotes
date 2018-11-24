
package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"` // We are using BSON to talk with the MongoDB database and JSON to talk to the client
	Title   string         `json:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey  string         `json:"apikey"`
}

// We switch on the HTTP method
func (s *Server) handlePolls(w http.ResponseWriter,
	r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlePollsGet(w, r)
		return
	case "POST":
		s.handlePollsPost(w, r)
		return
	case "DELETE":
		s.handlePollsDelete(w, r)
		return
	// A CORS browser will actually send a preflight request (with an HTTP method of OPTIONS)  asking for permission to make a DELETE request (listed in the Access-Control-Request-Method request header)
	case "OPTIONS":
		// the API will respond by setting the Access-Control-Allow-Methods header to DELETE, thus overriding the default * value that we set in our withCORS wrapper handler.
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}


// get poll data using id. return response
func (s *Server) handlePollsGet(w http.ResponseWriter, r *http.Request) {

	// create db copy
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")



// get poll data using id
func (s *Server) handlePollsGet(w http.ResponseWriter, r *http.Request) {
	//  create a copy of the database session that will allow us to interact with MongoDB
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	
	// We then build up an mgo.Query object by parsing the path.
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific poll
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all polls
		q = c.Find(nil)
	}
	var result []*poll
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}

// create a poll
func (s *Server) handlePollsPost(w http.ResponseWriter,
	r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	var p poll
	// decode response - should contain a representation of the poll object the client wants to create
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read poll from request", err)
		return
	}
	// generate a new unique ID for the poll and use the mgo package's Insert method to send it into the database
	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = apikey
	}
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError,
			"failed to insert poll", err)
		return
	}
	// set the Location header of the response and respond with a 201 http.StatusCreated message, pointing to the URL from which the newly created poll may be accessed.
	w.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

// remove poll from db
func (s *Server) handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	//  parse the path
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed,
			"Cannot delete all polls.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError,
			"failed to delete poll", err)
		return
	}
	// return a 200 Success response
	respond(w, r, http.StatusOK, nil)
}
