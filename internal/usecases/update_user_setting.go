// file: internal/usecases/update_user_setting.go
package usecases

import (
	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"
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
		return fmt.Errorf("invalid setting key: %d", setting)
	}
	return uc.repo.Patch(ctx, userID, setting, value)
}