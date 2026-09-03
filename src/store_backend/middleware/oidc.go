package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

const claimsKey contextKey = "claims"

type Claims struct {
	Subject         string `json:"sub"`
	AuthorizedParty string `json:"azp"`
	Scope           string `json:"scope"`
}

func OIDCVerification(verifier *oidc.IDTokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Fields(authHeader)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "missing or invalid bearer token", http.StatusUnauthorized)
				return
			}

			token, err := verifier.Verify(r.Context(), parts[1])
			if err != nil {
				http.Error(w, "token failed verification", http.StatusUnauthorized)
				return
			}

			var claims Claims

			if err := token.Claims(&claims); err != nil {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(), claimsKey, claims,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireScope(requiredScope string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(claimsKey).(Claims)
			if !ok {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			scopes := strings.Fields(claims.Scope)
			if !contains(scopes, requiredScope) {
				http.Error(w, fmt.Sprintf("lacking scope %s", requiredScope), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func contains(items []string, query string) bool {
	for _, s := range items {
		if s == query {
			return true
		}
	}
	return false
}
