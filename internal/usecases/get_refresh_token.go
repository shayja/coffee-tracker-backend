// file: internal/usecases/get_refresh_token.go
package usecases

import (
	"coffee-tracker-backend/internal/repositories"
	"context"
	"time"

	"github.com/google/uuid"
)

type GetRefreshTokenUseCase struct {
	authRepo repositories.AuthRepository
}

func NewGetRefreshTokenUseCase(authRepo repositories.AuthRepository) *GetRefreshTokenUseCase {
	return &GetRefreshTokenUseCase{authRepo: authRepo}
}

func (uc *GetRefreshTokenUseCase) Execute(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (string, time.Time, error) {
	return uc.authRepo.GetRefreshToken(ctx, userID, deviceID)
}