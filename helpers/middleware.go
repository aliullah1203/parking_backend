package helpers

import (
	"net/http"
)

const claimsContextKey = contextKey("claims")

// AuthorizeRole restricts access based on roles
func AuthorizeRole(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claimsInterface := r.Context().Value(claimsContextKey)
		if claimsInterface == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := claimsInterface.(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

// GetClaimsFromContext extracts claims from request context
