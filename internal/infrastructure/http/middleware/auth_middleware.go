// file: internal/infrastructure/http/middleware/auth_middleware.go
package middleware

import (
	"coffee-tracker-backend/internal/contextkeys"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and extracts the user ID
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	if jwtSecret == "" {
		panic("JWT secret must not be empty")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := extractAndValidateToken(r, jwtSecret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Attach userID to context
			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractAndValidateToken parses and validates the JWT token from the Authorization header
func extractAndValidateToken(r *http.Request, jwtSecret string) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return uuid.Nil, errors.New("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return uuid.Nil, errors.New("invalid authorization header format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token with claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, errors.New("token is not valid")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	// Optional: validate exp claim if present
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return uuid.Nil, errors.New("token expired")
		}
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return uuid.Nil, errors.New("user ID not found in token")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format in token")
	}

	return userID, nil
}
