package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateJWT creates a signed JWT token for testing or authentication
func GenerateJWT(secret string, userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),          // Supabase user ID
		"aud":  "authenticated",          // typical Supabase audience
		"role": "authenticated",          // user role
		"exp":  time.Now().Add(time.Hour).Unix(), // expires in 1 hour
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
