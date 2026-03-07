// file: internal/infrastructure/http/middleware/auth_middleware.go

package middleware

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"context"
	"net/http"
)

func AuthMiddleware(tokenService auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.ExtractBearerToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate the token and extract claims
			claims, err := tokenService.ValidateTokenString(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Reject refresh tokens — only access tokens are allowed on protected routes
			if tokenService.IsRefreshToken(claims) {
				http.Error(w, "refresh token cannot be used for API access", http.StatusUnauthorized)
				return
			}

			userID, err := tokenService.ExtractUserIDFromToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
