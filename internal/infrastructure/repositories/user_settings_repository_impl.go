package repositories

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserSettingsRepositoryImpl struct {
	db *sql.DB
}

func NewUserSettingsRepositoryImpl(db *sql.DB) repositories.UserSettingsRepository {
	return &UserSettingsRepositoryImpl{db: db}
}

func (r *UserSettingsRepositoryImpl) Set(ctx context.Context, userID uuid.UUID, key, value string) error {
	query := `
		INSERT INTO user_settings (user_id, key, value, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, key)
		DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query, userID, key, value, time.Now())
	return err
}

func (r *UserSettingsRepositoryImpl) Get(ctx context.Context, userID uuid.UUID, key string) (string, error) {
	var value string
	query := `SELECT value FROM user_settings WHERE user_id = $1 AND key = $2`
	err := r.db.QueryRowContext(ctx, query, userID, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *UserSettingsRepositoryImpl) GetAll(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	query := `SELECT key, value FROM user_settings WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *UserSettingsRepositoryImpl) Delete(ctx context.Context, userID uuid.UUID, key string) error {
	query := `DELETE FROM user_settings WHERE user_id = $1 AND key = $2`
	_, err := r.db.ExecContext(ctx, query, userID, key)
	return err
}
