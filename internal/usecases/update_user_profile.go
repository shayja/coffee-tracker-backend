// file: internal/usecases/update_user_profile.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/infrastructure/http/models"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

type UpdateUserProfileUseCase struct {
	userRepo repositories.UserRepository
}

func NewUpdateUserProfileUseCase(userRepo repositories.UserRepository) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{userRepo: userRepo}
}

func (uc *UpdateUserProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, req *models.UpdateUserProfileRequest) error {
	 if (req.Name == nil || *req.Name == "") && (req.Email == nil || *req.Email == "") {
        return ErrInvalidInput
    }

	err := uc.userRepo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return err
	}
	return nil
}
