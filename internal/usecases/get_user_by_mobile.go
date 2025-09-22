// file: internal/usecases/get_user_by_mobile.go
package usecases

import (
	"context"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"
)

// GetUserByMobileUseCase retrieves a user by mobile number
type GetUserByMobileUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetUserByMobileUseCase(userRepo repositories.UserRepository) *GetUserByMobileUseCase {
	return &GetUserByMobileUseCase{userRepo: userRepo}
}

func (uc *GetUserByMobileUseCase) Execute(ctx context.Context, mobile string) (*entities.User, error) {
	user, err := uc.userRepo.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err // or wrap with custom error like ErrUserNotFound
	}
	return user, nil
}
