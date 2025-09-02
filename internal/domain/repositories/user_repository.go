// file: internal/domain/repositories/user_repository.go
package repositories

import (
	"context"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
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
	UpdateAProfileImage(ctx context.Context, user *entities.User) error
	DeleteProfileImage(ctx context.Context, userID uuid.UUID, updatedAt time.Time) error

}
