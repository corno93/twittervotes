<<<<<<< HEAD
=======
// The http.Request object gives us access to every piece of information we might need about the underlying HTTP request
// However URL parsing is one such this http.Request does not solve well. While we can access a path (such as /people/1/books/2) as a string via the http.Request
// type's URL.Path field, there is no easy way to pull out the data encoded in the path, such as the people ID of1or the book ID of 2.

>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
package main

import (
	"strings"
)

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

<<<<<<< HEAD
// Pull out the data encoded in the path, such as the people ID of1or the book ID of 2 and return filled Path type
=======
// This simple parser provides a NewPath function that parses the specified path string and returns a new instance of the Path type.
>>>>>>> bb1838635a6c47ca84d1140e9b4b4da06f4b83bb
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
