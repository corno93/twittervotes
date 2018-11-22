package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":8081", "website address")
	flag.Parse()
	// uses the http.ServeMux type to serve static files from a folder called public.
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	log.Println("Serving website at:", *addr)

	http.ListenAndServe(*addr, mux)
}
