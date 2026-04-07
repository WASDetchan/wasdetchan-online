package util

import (
	"context"

	"github.com/WASDetchan/wasdetchan-online/auth"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func MakeContext(c *gin.Context) context.Context {
	return auth.MakeAuthContext(context.Background(), c)
}

func ServeComponentAt(link string, comp templ.Component, r *gin.Engine) {
	r.GET(link, func(c *gin.Context) {
		comp.Render(auth.MakeAuthContext(context.Background(), c), c.Writer)
	})
}
