// file: internal/usecases/save_refresh_token.go
package usecases

import (
	"coffee-tracker-backend/internal/repositories"
	"context"
	"time"

	"github.com/google/uuid"
)

type SaveRefreshTokenUseCase struct {
	authRepo repositories.AuthRepository
}

func NewSaveRefreshTokenUseCase(authRepo repositories.AuthRepository) *SaveRefreshTokenUseCase {
	return &SaveRefreshTokenUseCase{authRepo: authRepo}
}

func (uc *SaveRefreshTokenUseCase) Execute(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID, token string, expiresAt time.Time) error{
    return uc.authRepo.SaveRefreshToken(ctx, userID, deviceID, token, expiresAt)
}
