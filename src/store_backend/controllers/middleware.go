package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Middleware func(http.Handler) http.Handler

func IntegerPathParam(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			itemParam := chi.URLParam(r, key)
			var result int
			if parsed, err := strconv.Atoi(itemParam); err != nil {
				http.Error(w, "invalid path parameter", http.StatusBadRequest)
				return
			} else {
				result = parsed
			}
			ctx := context.WithValue(r.Context(), key, result)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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

func PassThru(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
