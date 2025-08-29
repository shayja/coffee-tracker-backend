// file: internal/domain/repositories/user_settings_repository.go
package repositories

import (
	"context"

	"github.com/google/uuid"
)

type UserSettingsRepository interface {
    Set(ctx context.Context, userID uuid.UUID, key, value string) error
    Get(ctx context.Context, userID uuid.UUID, key string) (string, error)
    GetAll(ctx context.Context, userID uuid.UUID) (map[string]string, error)
    Delete(ctx context.Context, userID uuid.UUID, key string) error
}
