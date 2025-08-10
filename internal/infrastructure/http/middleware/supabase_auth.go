package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// SupabaseAuthMiddleware validates Supabase JWT tokens
func SupabaseAuthMiddleware(supabaseKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// For now, use a default user ID when no auth is provided
				// In production, you should return 401 Unauthorized
				ctx := context.WithValue(r.Context(), "userID", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Extract token (Bearer token)
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// TODO: Validate JWT token with Supabase
			// For now, we'll use a hardcoded user ID
			// In production, decode the JWT and extract user ID
			userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
			
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
