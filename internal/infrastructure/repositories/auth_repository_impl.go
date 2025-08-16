package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// AuthRepositoryImpl implements AuthRepository with a SQL database.
type AuthRepositoryImpl struct {
	db *sql.DB
}

// NewAuthRepositoryImpl creates a new AuthRepositoryImpl
func NewAuthRepositoryImpl(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

// SaveRefreshToken inserts or updates a refresh token for a user
func (r *AuthRepositoryImpl) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt int64) error {
	query := `
		INSERT INTO user_refresh_tokens (user_id, token, expires_at, created_at, updated_at)
		VALUES ($1, $2, to_timestamp($3), NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

// GetRefreshToken retrieves the refresh token and expiry for a given user
func (r *AuthRepositoryImpl) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, int64, error) {
	var token string
	var expiresAt time.Time
	query := `SELECT token, expires_at FROM user_refresh_tokens WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&token, &expiresAt)
	if err != nil {
		return "", 0, err
	}
	return token, expiresAt.Unix(), nil
}

// DeleteRefreshToken removes a refresh token for a user
func (r *AuthRepositoryImpl) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
