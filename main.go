package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/feed.yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/xml")
		http.ServeFile(w, r, "/public/feed.yml")
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
