// file: internal/infrastructure/auth/jwt_service.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService implements TokenService
var _ TokenService = (*JWTService)(nil)

//
// ===========================
// 🧩 Implementation
// ===========================
//

// JWTService provides JWT generation and validation
type JWTService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	signingMethod jwt.SigningMethod
	nowFunc       func() time.Time // injected for testability
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		signingMethod: jwt.SigningMethodHS256,
		nowFunc:       func() time.Time { return time.Now().UTC() }, // can be replaced in tests
	}
}

//
// ===========================
// 🏷️ Token Generation
// ===========================
//

// GenerateAccessToken creates a short-lived JWT for API access.
func (s *JWTService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := s.nowFunc()
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"aud":  "authenticated",
		"role": "authenticated",
		"exp":  now.Add(s.accessExpiry).Unix(),
		"iat":  now.Unix(),
		"type": "access",
		"jti":  uuid.New().String(), // JWT ID for revocation
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateRefreshToken creates a long-lived JWT for session refresh.
func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	now := s.nowFunc()
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"aud":  "authenticated",
		"role": "refresh",
		"exp":  now.Add(s.refreshExpiry).Unix(),
		"iat":  now.Unix(),
		"type": "refresh",
		"jti":  uuid.New().String(),
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)
	return token.SignedString([]byte(s.secret))
}

//
// ===========================
// 🔍 Token Validation
// ===========================
//

// ValidateTokenString parses and validates a JWT string.
func (s *JWTService) ValidateTokenString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Check expiration defensively
	if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(s.nowFunc()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// ExtractUserIDFromToken retrieves the user ID from a validated token string.
func (s *JWTService) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateTokenString(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return uuid.Nil, ErrMissingUserID
	}

	return uuid.Parse(sub)
}

// TokenType returns the type of token ("access" or "refresh")
func (s *JWTService) TokenType(claims jwt.MapClaims) string {
	if t, ok := claims["type"].(string); ok {
		return t
	}
	return ""
}

// IsRefreshToken returns true if claims belong to a refresh token
func (s *JWTService) IsRefreshToken(claims jwt.MapClaims) bool {
	return s.TokenType(claims) == "refresh"
}

// ParseAndValidate is a helper that returns both userID and claims
func (s *JWTService) ParseAndValidate(tokenString string) (uuid.UUID, jwt.MapClaims, error) {
	claims, err := s.ValidateTokenString(tokenString)
	if err != nil {
		return uuid.Nil, nil, err
	}

	uid, err := s.ExtractUserIDFromToken(tokenString)
	return uid, claims, err
}

//
// ===========================
// ⏱️ Accessors
// ===========================
//

func (s *JWTService) AccessExpiry() time.Duration  { return s.accessExpiry }
func (s *JWTService) RefreshExpiry() time.Duration { return s.refreshExpiry }

//
// ===========================
// 🧪 Test helper
// ===========================
//

// WithNow allows tests to control current time.
func (s *JWTService) WithNow(f func() time.Time) *JWTService {
	s.nowFunc = f
	return s
}
