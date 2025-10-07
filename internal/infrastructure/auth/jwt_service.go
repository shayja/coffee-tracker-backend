// file: internal/infrastructure/auth/jwt_service.go
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(secret string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID.String(),
		"aud":   "authenticated",
		"role":  "authenticated",
		"exp":   time.Now().Add(s.accessExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"type":  "access",
		"token": uuid.New().String(), // Unique token ID for potential blacklisting
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID.String(),
		"exp":   time.Now().Add(s.refreshExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"type":  "refresh",
		"token": uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ValidateTokenString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Optional expiration check (already in claims["exp"])
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}

// ExtractUserIDFromToken gets the user ID from a validated token string.
func (s *JWTService) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateTokenString(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return uuid.Nil, errors.New("user ID not found in token")
	}

	return uuid.Parse(sub)
}

func (s *JWTService) IsRefreshToken(claims jwt.MapClaims) bool {
	tokenType, ok := claims["type"].(string)
	return ok && tokenType == "refresh"
}

func (s *JWTService) RefreshExpiry() time.Duration {
	return s.refreshExpiry
}

func (s *JWTService) AccessExpiry() time.Duration {
	return s.accessExpiry
}