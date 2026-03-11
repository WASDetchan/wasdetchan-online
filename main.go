package main

import (
	"context"
	"log"
	"net/http"

	"github.com/WASDetchan/wasdetchan-online/pages"
	"github.com/a-h/templ"
)

func main() {
	http.Handle("/home", templ.Handler(pages.Home()))

	home := pages.Home()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		home.Render(context.Background(), w)
	})

	http.HandleFunc("/feed.yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/xml")
		http.ServeFile(w, r, "/public/feed.yml")
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("/public"))))

	log.Fatal(http.ListenAndServe(":8082", nil))
}
