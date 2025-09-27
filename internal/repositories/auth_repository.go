// file: internal/repositories/auth_repository.go
package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuthRepository defines the contract for managing authentication-related data.
type AuthRepository interface {

	// SaveOTP saves a one-time password (OTP) for a user with an expiration time.
	SaveOTP(ctx context.Context, userID uuid.UUID, otp string, expiresAt time.Time) error

	// GetValidOTP retrieves a valid OTP for a user, checking if it has not expired.
	GetValidOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error)

	// InvalidateOTP marks an OTP as invalid, preventing its future use.
	InvalidateOTP(ctx context.Context, userID uuid.UUID, otp string) error

	// SaveRefreshToken inserts or replaces the refresh token for a user
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID, token string, expiresAt time.Time) error

	// GetRefreshToken retrieves the refresh token and expiry for a given user
	GetRefreshToken(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (token string, expiresAt time.Time, err error)

	// DeleteRefreshToken removes a refresh token for a user (e.g., logout)
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error

	// GetUserIDByRefreshToken retrieves the user ID associated with a given refresh token
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)

	// InvalidateAllUserTokens marks all users OTP as invalid,
	InvalidateAllUserTokens(ctx context.Context, userID uuid.UUID) error
}