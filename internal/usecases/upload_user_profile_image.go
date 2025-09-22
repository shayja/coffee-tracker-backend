// file: internal/usecases/upload_user_profile_image.go
package usecases

import (
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/storage"
	"coffee-tracker-backend/internal/infrastructure/utils"

	"github.com/google/uuid"
)


type UploadUserProfileImageUseCase struct {
    userRepo repositories.UserRepository
    storage storage.StorageService
    config  *config.Config
}

func NewUploadUserProfileImageUseCase(userRepo repositories.UserRepository, storage storage.StorageService, config  *config.Config) *UploadUserProfileImageUseCase {
	return &UploadUserProfileImageUseCase{ userRepo: userRepo, storage: storage, config: config }
}

func (uc *UploadUserProfileImageUseCase) Execute(ctx context.Context, userID uuid.UUID, filename string, file io.Reader) (string, error) {
   
    const avatarFileNameLangth = 10. 
    extension := path.Ext(filename)
	userFolderPath := fmt.Sprintf("%s/%s%s", userID, utils.GenerateString(avatarFileNameLangth), extension)

    // Upload to storage
    url, err := uc.storage.UploadFile(ctx, uc.config.ProfileImageBucket, userFolderPath, file, true)

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