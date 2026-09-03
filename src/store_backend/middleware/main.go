package middleware

import "net/http"

type contextKey string

type Middleware func(http.Handler) http.Handler

func PassThru(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
