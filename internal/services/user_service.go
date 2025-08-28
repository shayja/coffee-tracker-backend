// file: internal/services/user_service.go
package services

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByMobile(ctx context.Context, mobile string) (*entities.User, error) {

	user, err := s.userRepo.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return user, nil
}