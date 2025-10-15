package helpers

import (
	"context"
	"errors"
	"log"
	"net/http"
)

// Password helpers
func VerifyPassword(stored, provided string) bool {
	return stored == provided
}

func HashPassword(password string) string {
	return password
}

// AuthMiddleware validates token and stores claims in context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := ValidateToken(token)
		if err != nil {
			log.Println("Invalid token:", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Store claims in context
		ctx := contextWithClaims(r.Context(), claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// Dummy token validator
func ValidateToken(token string) (map[string]interface{}, error) {
	if token == "valid-token" {
		return map[string]interface{}{
			"user_id": "123",
			"role":    "ADMIN", // Add role here for testing
		}, nil
	}
	return nil, errors.New("invalid token")
}

// Context helpers
type claimsKeyType string

const claimsKey claimsKeyType = "claims"

func contextWithClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(claimsKey).(map[string]interface{})
	return claims, ok
}

// Role authorization middleware
func AuthorizeRole(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaimsFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
