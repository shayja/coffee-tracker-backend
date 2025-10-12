// file: internal/infrastructure/repositories/auth_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

// AuthRepositoryImpl implements repositories.AuthRepository using SQL database
type AuthRepositoryImpl struct {
	db *sql.DB
}

// NewAuthRepositoryImpl creates a new AuthRepositoryImpl
func NewAuthRepositoryImpl(db *sql.DB) repositories.AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

// SaveOTP inserts a new OTP record for a user
func (r *AuthRepositoryImpl) SaveOTP(ctx context.Context, userID uuid.UUID, otp string, expiresAt time.Time) error {
	query := `
		INSERT INTO user_otps (user_id, otp, expires_at, used)
		VALUES ($1, $2, $3, FALSE)
	`
	_, err := r.db.ExecContext(ctx, query, userID, otp, expiresAt)
	return err
}

// GetValidOTP checks if the OTP is valid and unused
func (r *AuthRepositoryImpl) GetValidOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error) {
	var valid bool

	query := `
		SELECT COUNT(*) > 0
		FROM user_otps
		WHERE user_id = $1 AND otp = $2 AND expires_at > $3 AND used = FALSE
	`
	err := r.db.QueryRowContext(ctx, query, userID, otp, utils.NowUTC()).Scan(&valid)
	if err != nil {
		return false, err
	}
	return valid, nil
}

// InvalidateOTP marks the OTP as used
func (r *AuthRepositoryImpl) InvalidateOTP(ctx context.Context, userID uuid.UUID, otp string) error {

	query := `UPDATE user_otps SET used = TRUE WHERE user_id = $1 AND otp = $2 AND expires_at > $3`
    _, err := r.db.ExecContext(ctx, query, userID, otp, utils.NowUTC())
    return err
}

// SaveRefreshToken inserts or updates a refresh token for a device
func (r *AuthRepositoryImpl) SaveRefreshToken(ctx context.Context, userID, deviceID uuid.UUID, token string, expiresAt time.Time) error {

	query := `
		INSERT INTO user_refresh_tokens (user_id, device_id, token, expires_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, device_id)
		DO UPDATE SET 
			token = EXCLUDED.token, 
			expires_at = EXCLUDED.expires_at,
			updated_at = $5
	`
	_, err := r.db.ExecContext(ctx, query, userID, deviceID, token, expiresAt, utils.NowUTC())
	return err
}

// GetRefreshToken retrieves a refresh token and its expiration
func (r *AuthRepositoryImpl) GetRefreshToken(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (string, time.Time, error) {
	var token string
	var expiresAt time.Time

	query := `
		SELECT token, expires_at
		FROM user_refresh_tokens
		WHERE user_id = $1 AND device_id = $2 AND expires_at > $3
	`
	err := r.db.QueryRowContext(ctx, query, userID, deviceID, utils.NowUTC()).Scan(&token, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", time.Time{}, repositories.ErrNotFound
		}
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

// DeleteRefreshToken removes a refresh token for a device
func (r *AuthRepositoryImpl) DeleteRefreshToken(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) error {
	query := `
		DELETE FROM user_refresh_tokens
		WHERE user_id = $1 AND device_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, deviceID)
	return err
}

// GetUserIDByRefreshToken returns the user ID associated with a refresh token
func (r *AuthRepositoryImpl) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := `
		SELECT user_id
		FROM user_refresh_tokens
		WHERE token = $1 AND expires_at > $2
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, refreshToken, utils.NowUTC()).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, repositories.ErrNotFound
		}
		return uuid.Nil, err
	}

	return userID, nil
}

// InvalidateAllUserTokens deletes all refresh tokens for a user
func (r *AuthRepositoryImpl) InvalidateAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
