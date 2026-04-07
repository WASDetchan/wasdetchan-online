//go:generate go tool templ generate
//go:generate go tool sqlc generate

package main

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/WASDetchan/wasdetchan-online/auth"
	"github.com/WASDetchan/wasdetchan-online/pages"
	"github.com/WASDetchan/wasdetchan-online/pages/articles"
	"github.com/WASDetchan/wasdetchan-online/repository"
	"github.com/WASDetchan/wasdetchan-online/util"
)

func main() {
	queries, err := repository.InitPostgres()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return
	}

	key := make([]byte, 64)
	rand.Read(key)
	store := cookie.NewStore(key)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(sessions.Sessions("session", store))

	auth.RegisterAuth(r, queries)

	util.ServeComponentAt("/home", pages.Home(), r)
	util.ServeComponentAt("/", pages.Home(), r)
	util.ServeComponentAt("/articles", pages.Articles(), r)

	articles.HelloWorldInfo().Register(articles.HelloWorld(), r)
	articles.ThisWebsiteInfo().Register(articles.ThisWebsite(), r)

	r.GET("/feed.yml", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/xml")
		w := ctx.Writer
		r := ctx.Request
		http.ServeFile(w, r, "/public/feed.yml")
	})

	r.Static("/public/", "/public/")
	r.Static("/static/", "/static/")

	r.Run(":8082")

}
