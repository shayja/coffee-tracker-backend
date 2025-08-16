package services

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByMobile(ctx context.Context, mobile string) (*entities.User, error) {

	user, err := s.userRepo.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return user, nil
}