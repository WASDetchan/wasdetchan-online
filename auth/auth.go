package auth

import (
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

type AuthInfo struct {
	Email    string
	Username string
}

type UserKey struct{}

func RegisterAuth(r *gin.Engine) {
	gob.Register(AuthInfo{})

	key := make([]byte, 64)
	rand.Read(key)
	gothicStore := gsessions.NewCookieStore(key)
	gothicStore.MaxAge(86400 * 30)
	gothicStore.Options.Path = "/"
	gothicStore.Options.HttpOnly = true
	gothicStore.Options.Secure = true
	gothic.Store = gothicStore

	var providerNames []string
	var providers []goth.Provider

	if os.Getenv("GITHUB_KEY") != "" && os.Getenv("GITHUB_SECRET") != "" {
		providerNames = append(providerNames, "github")
		providers = append(providers,
			github.New(
				os.Getenv("GITHUB_KEY"),
				os.Getenv("GITHUB_SECRET"),
				fmt.Sprintf("%v/auth/github/callback", os.Getenv("URL")),
			),
		)
	}

	goth.UseProviders(providers...)

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		query := c.Request.URL.Query()
		query.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = query.Encode()

		req := c.Request
		res := c.Writer

		gothUser, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		session := sessions.Default(c)

		user := AuthInfo{
			gothUser.Email,
			gothUser.Name,
		}

		session.Set(AuthInfo{}, user)
		if err := session.Save(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		log.Printf("Logged in as %v with %v", user, c.Param("provider"))

		c.Redirect(http.StatusTemporaryRedirect, "/") // TODO: redirect to source
	})

	r.GET("/auth/:provider", func(c *gin.Context) {
		log.Print(c.Param("provider"))

		query := c.Request.URL.Query()
		query.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = query.Encode()

		req := c.Request
		res := c.Writer

		if user, err := gothic.CompleteUserAuth(res, req); err == nil {

			session := sessions.Default(c)

			user := AuthInfo{
				user.Email,
				user.Name,
			}

			session.Set(UserKey{}, user)
			session.Save()

			log.Printf("Logged in as %v", user)

			c.Redirect(http.StatusTemporaryRedirect, "/") // TODO: redirect to source
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})
}
