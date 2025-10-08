// file: internal/usecases/generate_otp.go
package usecases

import (
	"coffee-tracker-backend/internal/infrastructure/config"
	"coffee-tracker-backend/internal/infrastructure/notifications"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/repositories"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GenerateOtpUseCase struct {
	authRepo repositories.AuthRepository
	config  *config.Config
	smsService notifications.SMSService
}

func NewGenerateOtpUseCase(authRepo repositories.AuthRepository, config  *config.Config, smsService notifications.SMSService) *GenerateOtpUseCase {
	return &GenerateOtpUseCase{authRepo: authRepo, config: config, smsService: smsService }
}

func (uc *GenerateOtpUseCase) Execute(ctx context.Context, userID uuid.UUID, mobile string) (error) {
	// generate random N-digit OTP
	strength := config.OtpStrength(uc.config.OtpStrength) // cast string → OtpStrength
	otp, err := utils.GenerateOTP(strength)
	if err != nil {
		return err
	}

	// OTP valid for N minutes
	expiresAt := time.Now().Add(5 * time.Minute).UTC()

	// save OTP to DB
	err = uc.authRepo.SaveOTP(ctx, userID, otp, expiresAt)
	if err != nil {
		return err
	}

	// Send SMS here
	if err := uc.smsService.SendOTP(userID, mobile, otp); err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	
	return nil
}