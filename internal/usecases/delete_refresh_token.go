// file: internal/usecases/delete_refresh_token.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"context"

	"github.com/google/uuid"
)

type DeleteRefreshTokenUseCase struct {
	authRepo repositories.AuthRepository
}

func NewDeleteRefreshTokenUseCase(authRepo repositories.AuthRepository) *DeleteRefreshTokenUseCase {
	return &DeleteRefreshTokenUseCase{authRepo: authRepo}
}

func (uc *DeleteRefreshTokenUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	return uc.authRepo.DeleteRefreshToken(ctx, userID)
}