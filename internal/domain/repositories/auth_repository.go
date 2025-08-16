package repositories

import (
	"context"

	"github.com/google/uuid"
)

// AuthRepository defines the contract for managing authentication-related data.
type AuthRepository interface {
	// SaveRefreshToken inserts or replaces the refresh token for a user
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt int64) error

	// GetRefreshToken retrieves the refresh token and expiry for a given user
	GetRefreshToken(ctx context.Context, userID uuid.UUID) (token string, expiresAt int64, err error)

	// DeleteRefreshToken removes a refresh token for a user (e.g., logout)
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error

	// GetUserIDByRefreshToken retrieves the user ID associated with a given refresh token
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
}