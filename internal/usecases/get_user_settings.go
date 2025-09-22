// file: internal/usecases/get_user_settings.go
package usecases

import (
	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"
	"context"

	"github.com/google/uuid"
)

type GetUserSettingsUseCase struct {
	settingsRepo repositories.UserSettingsRepository
}

func NewGetUserSettingsUseCase(settingsRepo repositories.UserSettingsRepository) *GetUserSettingsUseCase {
	return &GetUserSettingsUseCase{settingsRepo: settingsRepo}
}

func (uc *GetUserSettingsUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entities.UserSettings, error) {
	return uc.settingsRepo.Get(ctx, userID)
}
