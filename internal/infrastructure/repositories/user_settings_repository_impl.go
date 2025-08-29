// file: internal/infrastructure/repositories/user_settings_repository_impl.go
package repositories

import (
	"coffee-tracker-backend/internal/domain/entities"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserSettingsRepositoryImpl struct {
    db *sql.DB
}

func NewUserSettingsRepositoryImpl(db *sql.DB) *UserSettingsRepositoryImpl {
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
func (r *UserSettingsRepositoryImpl) Patch(ctx context.Context, userID uuid.UUID, updates map[entities.Setting]interface{}) error {
    if len(updates) == 0 {
        return nil
    }

    setClauses := []string{}
    args := []interface{}{}
    i := 1

    for setting, value := range updates {
        if !setting.IsValid() {
            return fmt.Errorf("invalid setting: %s", setting)
        }
        setClauses = append(setClauses, fmt.Sprintf("%s = $%d", setting, i))
        args = append(args, value)
        i++
    }

    query := fmt.Sprintf(
        `UPDATE user_settings SET %s, updated_at = now() WHERE user_id = $%d`,
        strings.Join(setClauses, ", "),
        i,
    )
    args = append(args, userID)

    _, err := r.db.ExecContext(ctx, query, args...)
    return err
}

// Reset sets a specific setting to its default (e.g. false)
func (r *UserSettingsRepositoryImpl) Reset(ctx context.Context, userID uuid.UUID, setting entities.Setting) error {
    if !setting.IsValid() {
        return fmt.Errorf("invalid setting: %s", setting)
    }

    query := fmt.Sprintf(`UPDATE user_settings SET %s = false, updated_at = now() WHERE user_id = $1`, setting)
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}
