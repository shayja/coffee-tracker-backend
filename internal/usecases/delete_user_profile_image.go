// file: internal/usecases/delete_user_profile_image.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type DeleteUserProfileImageUseCase struct {
	userRepo repositories.UserRepository
}

func NewDeleteUserProfileImageUseCase(userRepo repositories.UserRepository) *DeleteUserProfileImageUseCase {
	return &DeleteUserProfileImageUseCase{userRepo: userRepo}
}

func (uc *DeleteUserProfileImageUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.DeleteProfileImage(ctx, userID)
}
