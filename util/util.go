package util

import (
	"github.com/WASDetchan/wasdetchan-online/core"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func ServeComponentAt(link string, comp templ.Component, r *gin.RouterGroup) {
	r.GET(link, func(c *gin.Context) {
		comp.Render(core.GetContext(c), c.Writer)
	})
}
