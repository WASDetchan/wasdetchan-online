//go:generate go tool templ generate
//go:generate go tool sqlc generate

package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/WASDetchan/wasdetchan-online/auth"
	"github.com/WASDetchan/wasdetchan-online/pages"
	"github.com/WASDetchan/wasdetchan-online/repository"
	"github.com/a-h/templ"
)

func makeContext(c *gin.Context) context.Context {
	return context.WithValue(context.Background(), auth.UserKey{}, sessions.Default(c).Get(auth.UserKey{}))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

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

	home := templ.Handler(pages.Home())
	r.GET("/home", func(c *gin.Context) {
		home.Component.Render(makeContext(c), c.Writer)
	})
	r.GET("/", func(c *gin.Context) {
		home.Component.Render(makeContext(c), c.Writer)
	})

	article := templ.Handler(pages.HelloWorld())
	r.GET("/articles/hello_world", func(c *gin.Context) {
		article.Component.Render(makeContext(c), c.Writer)
	})

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
