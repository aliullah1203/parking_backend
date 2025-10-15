package helpers

import (
	"context"
	"net/http"
	"strings"
)

// Key type for storing claims in context
type contextKey string

const claimsContextKey = contextKey("claims")

// AuthMiddleware validates JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// store claims in request context
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
func GetClaimsFromContext(r *http.Request) *Claims {
	claimsInterface := r.Context().Value(claimsContextKey)
	if claimsInterface == nil {
		return nil
	}
	if claims, ok := claimsInterface.(*Claims); ok {
		return claims
	}
	return nil
}
