package services

import (
	"coffee-tracker-backend/internal/infrastructure/repositories"
	"context"

	"github.com/google/uuid"
)

type UserSettingsService struct {
	settingsRepo *repositories.UserSettingsRepositoryImpl
}

func NewUserSettingsService(settingsRepo *repositories.UserSettingsRepositoryImpl) *UserSettingsService {
	return &UserSettingsService{settingsRepo: settingsRepo}
}

func (s *UserSettingsService) Set(ctx context.Context, userID uuid.UUID, key, value string) error {
	return s.settingsRepo.Set(ctx, userID, key, value)
}

func (s *UserSettingsService) Get(ctx context.Context, userID uuid.UUID, key string) (string, error) {
	return s.settingsRepo.Get(ctx, userID, key)
}

func (s *UserSettingsService) GetAll(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	return s.settingsRepo.GetAll(ctx, userID)
}

func (s *UserSettingsService) Delete(ctx context.Context, userID uuid.UUID, key string) error {
	return s.settingsRepo.Delete(ctx, userID, key)
}
