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
				w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_request\"")
				http.Error(w, "missing or malformed bearer token", http.StatusBadRequest)
				return
			}

			token, err := verifier.Verify(r.Context(), parts[1])
			if err != nil {
				w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
				http.Error(w, "token failed verification", http.StatusUnauthorized)
				return
			}

			var claims Claims

			if err := token.Claims(&claims); err != nil {
				w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
				http.Error(w, "invalid or malformed claims", http.StatusUnauthorized)
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
				w.Header().Set(
					"WWW-Authenticate",
					fmt.Sprintf(`Bearer error="insufficient_scope", scope="%s"`, requiredScope))
				http.Error(w, fmt.Sprintf("lacking scope %s", requiredScope), http.StatusForbidden)
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
