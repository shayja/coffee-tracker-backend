// file: internal/usecases/update_user_setting.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UpdateUserSettingUseCase struct {
    repo repositories.UserSettingsRepository
}

func NewUpdateUserSettingUseCase(repo repositories.UserSettingsRepository) *UpdateUserSettingUseCase {
    return &UpdateUserSettingUseCase{repo: repo}
}

// Execute updates a single user setting
func (uc *UpdateUserSettingUseCase) Execute(ctx context.Context, userID uuid.UUID, setting entities.Setting, value interface{}) error {
    if !setting.IsValid() {
        return fmt.Errorf("invalid setting key: %s", setting)
    }

    updates := map[entities.Setting]interface{}{
        setting: value,
    }

    return uc.repo.Patch(ctx, userID, updates)
}
