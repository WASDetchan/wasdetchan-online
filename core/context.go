package core

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

type ContextKey struct{}

func PushContext(c *gin.Context, key any, value any) {
	log.Printf("Pushing %v: %v into context", key, value)
	ctxi, _ := c.Get(ContextKey{})
	ctx, _ := ctxi.(context.Context)
	if ctx == nil {
		log.Print("Creating new context")
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, key, value)

	c.Set(ContextKey{}, ctx)

	GetContext(c)
}

func GetContext(c *gin.Context) (ctx context.Context) {
	log.Print("Getting context")
	ctxi, _ := c.Get(ContextKey{})
	ctx, _ = ctxi.(context.Context)
	if ctx == nil {
		log.Print("Creating new context")
		ctx = context.Background()
	}
	return
}
