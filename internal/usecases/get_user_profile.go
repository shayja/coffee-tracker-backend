package usecases

import (
	"context"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetUserProfileUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetUserProfileUseCase(userRepo repositories.UserRepository) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{userRepo: userRepo}
}

func (uc *GetUserProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}
