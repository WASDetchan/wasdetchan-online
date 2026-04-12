package auth

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/WASDetchan/wasdetchan-online/core"
	"github.com/WASDetchan/wasdetchan-online/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5"

	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

type AuthInfo struct {
	User         *repository.User
	Capabilities Capabilities
}

func AssertAuth(c *gin.Context) *repository.User {
	info := GetAuthInfo(c)
	if info.User == nil {
		log.Panic("authenticated assert failed")
	}
	return info.User
}

func GetAuthInfo(c *gin.Context) AuthInfo {
	info, _ := c.Get(AuthKey{})
	return info.(AuthInfo)
}

type AuthKey struct{}

type RedirectKey struct{}

func RedirectToAuth(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set(RedirectKey{}, c.Request.URL.String())
	if err := sess.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/auth")
}

func EnsureAuthorized(capabilities ...Capability) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Printf("authorizing user...")
		user, _ := c.Get(AuthKey{})
		info, authenticated := user.(AuthInfo)
		if !authenticated || info.User == nil {
			RedirectToAuth(c)
			c.Abort()
			return
		}
		for _, capability := range capabilities {
			if !info.Capabilities.Has(capability) {
				c.String(http.StatusForbidden, "")
				c.Abort()
				return
			}
		}
		log.Printf("user is authorized %v", info)

		c.Next()
	}
}

func RegisterAuth(r *gin.Engine, q *repository.Queries) {
	gob.Register(AuthInfo{})
	gob.Register(AuthKey{})
	gob.Register(RedirectKey{})

	r.Use(middleware)

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

	r.GET("/auth", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/auth/github") })

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		if err := complpeteAuth(c, q); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	r.GET("/auth/:provider", func(c *gin.Context) {
		if err := complpeteAuth(c, q); err != nil {
			beginAuth(c)
			return
		}
	})
}

func middleware(c *gin.Context) {
	info, loggedIn := sessions.Default(c).Get(AuthKey{}).(AuthInfo)
	if !loggedIn {
		token := c.Request.Header.Get("Authorization")
		if token != "" {
			info = authToken(token, repository.GetQueries(c))
			log.Print(info)
		}

	}

	core.PushContext(c, AuthKey{}, info)
	c.Set(AuthKey{}, info)
}

func getUser(q *repository.Queries, gothUser goth.User) (repository.User, bool, error) {
	ctx := context.Background()
	if user, err := q.GetUserWithEmail(ctx, gothUser.Email); err == nil {
		return user, false, nil
	} else {
		if err != pgx.ErrNoRows {
			return repository.User{}, false, fmt.Errorf("error getting user: %v", err)
		}
		user, err = q.CreateUser(ctx,
			repository.CreateUserParams{
				Name:    gothUser.Name,
				Email:   gothUser.Email,
				IsAdmin: IsOwner(gothUser.Email),
			},
		)
		if err != nil {
			return repository.User{}, false, fmt.Errorf("error creating user: %v", err)
		}
		return user, true, nil
	}
}

func complpeteAuth(c *gin.Context, q *repository.Queries) error {
	query := c.Request.URL.Query()
	query.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = query.Encode()

	req := c.Request
	res := c.Writer

	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return err
	}

	session := sessions.Default(c)

	user, created, err := getUser(q, gothUser)
	if err != nil {
		return err
	}

	info := AuthInfo{&user, AuthenticatedCapabilities()}

	session.Set(AuthKey{}, info)
	if err := session.Save(); err != nil {
		return err
	}

	if created {
		log.Printf("Signed up as %v with %v", user, c.Param("provider"))
	} else {
		log.Printf("Logged in as %v with %v", user, c.Param("provider"))
	}

	redirect, succ := session.Get(RedirectKey{}).(string)

	if !succ || redirect == "" {
		redirect = "/"
	}

	c.Redirect(http.StatusTemporaryRedirect, redirect)
	return nil
}

func beginAuth(c *gin.Context) {
	query := c.Request.URL.Query()
	query.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = query.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}
