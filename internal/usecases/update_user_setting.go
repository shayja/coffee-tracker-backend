// file: internal/usecases/update_user_setting.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"context"

	"github.com/google/uuid"
)

type UpdateUserSettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateUserSettingUseCase struct {
	settingsRepo repositories.UserSettingsRepository
}

func NewUpdateUserSettingUseCase(settingsRepo repositories.UserSettingsRepository) *UpdateUserSettingUseCase {
	return &UpdateUserSettingUseCase{settingsRepo: settingsRepo}
}

func (uc *UpdateUserSettingUseCase) Execute(ctx context.Context, userID uuid.UUID, req UpdateUserSettingRequest) error {
	return uc.settingsRepo.Set(ctx, userID, req.Key, req.Value)
}
