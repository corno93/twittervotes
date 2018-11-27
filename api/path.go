package main

import (
	"strings"
)

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

// Pull out the data encoded in the path, such as the people ID and return filled Path type
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
