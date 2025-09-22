// file: internal/usecases/get_user_by_id.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

	"github.com/google/uuid"
)

// GetUserByIDUseCase retrieves a user by ID
type GetUserByIDUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetUserByIDUseCase(userRepo repositories.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{userRepo: userRepo}
}

func (uc *GetUserByIDUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err // or wrap with custom error like ErrUserNotFound
	}
	return user, nil
}
