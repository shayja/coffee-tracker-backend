// file: internal/infrastructure/repositories/user_settings_repository_impl.go
package repositories

import (
	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type UserSettingsRepositoryImpl struct {
    db *sql.DB
}

func NewUserSettingsRepositoryImpl(db *sql.DB) repositories.UserSettingsRepository {
    return &UserSettingsRepositoryImpl{db: db}
}

// Get returns all user settings as a struct
func (r *UserSettingsRepositoryImpl) Get(ctx context.Context, userID uuid.UUID) (*entities.UserSettings, error) {
    query := `
        SELECT user_id, biometric_enabled, dark_mode, notifications_enabled, created_at, updated_at
        FROM user_settings
        WHERE user_id = $1
    `
    row := r.db.QueryRowContext(ctx, query, userID)

    var s entities.UserSettings
    if err := row.Scan(&s.UserID, &s.BiometricEnabled, &s.DarkMode, &s.NotificationsEnabled, &s.CreatedAt, &s.UpdatedAt); err != nil {
        return nil, err
    }

    return &s, nil
}

// Patch updates one or more user settings dynamically
func (r *UserSettingsRepositoryImpl) Patch(ctx context.Context, userID uuid.UUID, setting entities.Setting, value interface{}) error {
	column := setting.ColumnName()
	if column == "" {
		return fmt.Errorf("unknown setting: %d", setting)
	}

	query := fmt.Sprintf(`UPDATE user_settings SET %s = $1, updated_at = NOW() WHERE user_id = $2`, column)
	_, err := r.db.ExecContext(ctx, query, value, userID)
	return err
}


// Reset sets a specific setting to its default (e.g. false)
func (r *UserSettingsRepositoryImpl) Reset(ctx context.Context, userID uuid.UUID, setting entities.Setting) error {
    column := setting.ColumnName()
	if column == "" {
		return fmt.Errorf("unknown setting: %d", setting)
	}

    query := fmt.Sprintf(`UPDATE user_settings SET %s = false, updated_at = now() WHERE user_id = $1`, column)
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}
