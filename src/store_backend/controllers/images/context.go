package images

import "context"

func ItemIdFromContext(ctx context.Context) int {
	val, ok := ctx.Value("itemId").(int)
	if !ok {
		panic("expected itemId in context")
	}
	return val
}

func ImageIdFromContext(ctx context.Context) int {
	val, ok := ctx.Value("imageId").(int)
	if !ok {
		panic("expected imageId in context")
	}
	return val
}
