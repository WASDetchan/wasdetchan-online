package core

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ContextKey struct{}

func PushContext(c *gin.Context, key any, value any) {
	ctxi, _ := c.Get(ContextKey{})
	ctx, _ := ctxi.(context.Context)
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, key, value)

	c.Set(ContextKey{}, ctx)

	GetContext(c)
}

func GetContext(c *gin.Context) (ctx context.Context) {
	ctxi, _ := c.Get(ContextKey{})
	ctx, _ = ctxi.(context.Context)
	if ctx == nil {
		ctx = context.Background()
	}
	return
}
