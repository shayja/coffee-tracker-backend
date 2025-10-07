// file: internal/usecases/validate_otp.go
package usecases

import (
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/repositories"
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
	valid, err := uc.authRepo.GetValidOTP(ctx, userID, otp)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, ErrInvalidOTP // define this in `usecases/errors.go`
	}
	// Invalidate OTP after use
	if err := uc.authRepo.InvalidateOTP(ctx, userID, otp); err != nil {
		return true, nil
	}
	return true, nil

	
}
