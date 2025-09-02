package usecases

import (
	"context"
	"io"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/storage"

	"github.com/google/uuid"
)

type UploadUserProfileImageResult struct {
	AvatarURL string `json:"avatar_url"`
}

type UploadUserProfileImageUseCase struct {
    storage storage.StorageService
    userRepo repositories.UserRepository
}


func NewUploadUserProfileImageUseCase(userRepo repositories.UserRepository) *UploadUserProfileImageUseCase {
	return &UploadUserProfileImageUseCase{userRepo: userRepo}
}


func (uc *UploadUserProfileImageUseCase) Execute(ctx context.Context, userID uuid.UUID, filename string, file io.Reader) (string, error) {
    // Upload to storage (Supabase)
    url, err := uc.storage.UploadFile(ctx, "avatars", filename, file)
    if err != nil {
        return "", err
    }

    // Update user's avatar_url in DB
    user := &entities.User{
        ID:        userID,
        AvatarURL: url,
        UpdatedAt: time.Now().UTC(),
    }
    if err := uc.userRepo.UpdateAProfileImage(ctx, user); err != nil {
        return "", err
    }

    return url, nil
}