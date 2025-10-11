// file: internal/infrastructure/auth/auth/errors.go
package auth

import "errors"

// Standard JWT errors
var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidClaims = errors.New("invalid token claims")
	ErrMissingUserID = errors.New("user ID not found in token")
)
