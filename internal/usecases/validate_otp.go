// file: internal/usecases/validate_otp.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/config"
	"context"

	"github.com/google/uuid"
)

type ValidateOtpUseCase struct {
	authRepo repositories.AuthRepository
	config  *config.Config
}

func NewValidateOtpUseCase(authRepo repositories.AuthRepository, config  *config.Config) *ValidateOtpUseCase {
	return &ValidateOtpUseCase{authRepo: authRepo, config: config}
}

func (uc *ValidateOtpUseCase) Execute(ctx context.Context, userID uuid.UUID, otp string) (bool, error) {
	
	if otp == uc.config.MagicOtp {
		return true, nil
	}
	// Check if OTP is valid and not expired
	valid, err := uc.authRepo.GetValidOTP(ctx, userID.String(), otp)
	if err != nil || !valid {
		return false, err
	}

	// Invalidate OTP after use
	_ = uc.authRepo.InvalidateOTP(ctx, userID.String(), otp)

	return true, nil
}
