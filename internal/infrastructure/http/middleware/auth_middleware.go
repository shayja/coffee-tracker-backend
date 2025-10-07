// file: internal/infrastructure/http/middleware/auth_middleware.go

package middleware

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"context"
	"net/http"
)

func AuthMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.ExtractBearerToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := jwtService.ExtractUserIDFromToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
