// file: internal/infrastructure/repositories/auth_repository_impl.go
package repositories

import (
	"coffee-tracker-backend/internal/domain/repositories"
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
	func NewAuthRepositoryImpl(db *sql.DB) repositories.AuthRepository {
		return &AuthRepositoryImpl{db: db}
	}

	func (r *AuthRepositoryImpl) SaveOTP(ctx context.Context, userID uuid.UUID, otp string, expiresAt time.Time) error {
		query := `
			INSERT INTO user_otps (user_id, otp, expires_at, used)
			VALUES ($1, $2, $3, FALSE)
		`
		_, err := r.db.ExecContext(ctx, query, userID, otp, expiresAt)
		return err
	}
	

	func (r *AuthRepositoryImpl) GetValidOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error) {
		var valid bool
		query := `
			SELECT COUNT(*) > 0
			FROM user_otps
			WHERE user_id = $1 AND otp = $2 AND expires_at > NOW() AND used = FALSE
		`
		err := r.db.QueryRowContext(ctx, query, userID, otp).Scan(&valid)
		return valid, err
	}

	func (r *AuthRepositoryImpl) InvalidateOTP(ctx context.Context, userID uuid.UUID, otp string) error {
		query := `UPDATE user_otps SET used = TRUE WHERE user_id = $1 AND otp = $2`
		_, err := r.db.ExecContext(ctx, query, userID, otp)
		return err
	}

	func (r *AuthRepositoryImpl) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
		query := `
			INSERT INTO user_refresh_tokens (user_id, token, expires_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) 
			DO UPDATE SET 
				token = EXCLUDED.token, 
				expires_at = EXCLUDED.expires_at,
				updated_at = NOW()
		`
		_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
		return err
	}

	func (r *AuthRepositoryImpl) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
		var token string
		var expiresAt time.Time
		query := `SELECT token, expires_at FROM user_refresh_tokens WHERE user_id = $1 AND expires_at > NOW()`
		err := r.db.QueryRowContext(ctx, query, userID).Scan(&token, &expiresAt)
		if err != nil {
			return "", time.Time{}, err
		}
		return token, expiresAt, nil
	}

	func (r *AuthRepositoryImpl) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
		query := `DELETE FROM user_refresh_tokens WHERE user_id = $1`
		_, err := r.db.ExecContext(ctx, query, userID)
		return err
	}

	func (r *AuthRepositoryImpl) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
		var userID uuid.UUID
		query := `
			SELECT user_id
			FROM user_refresh_tokens
			WHERE token = $1 AND expires_at > NOW()
			LIMIT 1
		`
		err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return uuid.Nil, repositories.ErrNotFound
			}
			return uuid.Nil, err
		}
		return userID, nil
	}

	func (r *AuthRepositoryImpl) InvalidateAllUserTokens(ctx context.Context, userID uuid.UUID) error {
		query := `DELETE FROM user_refresh_tokens WHERE user_id = $1`
		_, err := r.db.ExecContext(ctx, query, userID)
		return err
	}
