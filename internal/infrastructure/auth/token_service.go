package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenService defines the contract for any token-based authentication system.
type TokenService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateTokenString(tokenString string) (jwt.MapClaims, error)
	ExtractUserIDFromToken(tokenString string) (uuid.UUID, error)
	IsRefreshToken(claims jwt.MapClaims) bool
	AccessExpiry() time.Duration
	RefreshExpiry() time.Duration
}