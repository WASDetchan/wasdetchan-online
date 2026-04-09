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
	"github.com/WASDetchan/wasdetchan-online/receipt"
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
	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(sessions.Sessions("session", store))
	r.Use(repository.Middleware(queries))

	auth.RegisterAuth(r, queries)
	auth.MakeOwnerAdmin(queries)

	g := r.Group("/")

	util.ServeComponentAt("/home", pages.Home(), g)
	util.ServeComponentAt("/", pages.Home(), g)
	util.ServeComponentAt("/articles", pages.Articles(), g)

	articles.HelloWorldInfo().Register(articles.HelloWorld(), g)
	articles.ThisWebsiteInfo().Register(articles.ThisWebsite(), g)

	r.GET("/feed.yml", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/xml")
		w := ctx.Writer
		r := ctx.Request
		http.ServeFile(w, r, "/public/feed.yml")
	})

	authenticated := g.Group("/", auth.EnsureAuthenticated)
	util.ServeComponentAt("/receipts", pages.Receipts(), authenticated)
	authenticated.POST("/receipts", receipt.HandlePostReceipt)

	r.Static("/public/", "/public/")
	r.Static("/static/", "/static/")

	r.Run(":8082")

}
