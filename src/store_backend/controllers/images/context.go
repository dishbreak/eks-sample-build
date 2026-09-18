package images

import (
	"context"

	"github.com/dishbreak/sample-store-backend/middleware"
)

func ItemIdFromContext(ctx context.Context) int {
	val, ok := ctx.Value(middleware.ContextKey("itemId")).(int)
	if !ok {
		panic("expected itemId in context")
	}
	return val
}

func ImageIdFromContext(ctx context.Context) int {
	val, ok := ctx.Value(middleware.ContextKey("imageId")).(int)
	if !ok {
		panic("expected imageId in context")
	}
	return val
}
