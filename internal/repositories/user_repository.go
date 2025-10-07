// file: internal/repositories/user_repository.go
package repositories

import (
	"context"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/infrastructure/http/dto"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByMobile(ctx context.Context, mobile string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserProfileRequest) error
	UpdateProfileImage(ctx context.Context, user *entities.User) error
	DeleteProfileImage(ctx context.Context, userID uuid.UUID) error

}
