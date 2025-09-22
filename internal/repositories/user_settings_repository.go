// file: internal/repositories/user_settings_repository.go
package repositories

import (
	"coffee-tracker-backend/internal/entities"
	"context"

	"github.com/google/uuid"
)

type UserSettingsRepository interface {
    Get(ctx context.Context, userID uuid.UUID) (*entities.UserSettings, error)
    Patch(ctx context.Context, userID uuid.UUID, setting entities.Setting, value interface{}) error
    Reset(ctx context.Context, userID uuid.UUID, setting entities.Setting) error
}