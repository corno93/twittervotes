package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// These two functions abstract the decoding and encoding of data from and to the Request and ResponseWriter objects, respectively
// if we decide to add support for other representations or switch to a binary protocol instead, we only need to touch these two functions.
func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// makes life easier to write the status code and some data to the ResponseWriter object using our encodeBody helper.
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, r, data)
	}
}

// gives us an interface similar to the respond function, but the data written will be enveloped in an error object in order to make it clear that something went wrong.
func respondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

// add an HTTP-error-specific helper that will generate the correct message for us using the http.StatusText function from the Go standard library:
func respondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
	respondErr(w, r, status, http.StatusText(status))
}
