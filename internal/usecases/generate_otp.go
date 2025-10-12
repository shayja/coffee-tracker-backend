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
	smsService notifications.SMSService
	strength  config.OtpStrength
}

func NewGenerateOtpUseCase(authRepo repositories.AuthRepository, smsService notifications.SMSService, strength config.OtpStrength) *GenerateOtpUseCase {
	return &GenerateOtpUseCase{authRepo: authRepo, smsService: smsService, strength: strength }
}

func (uc *GenerateOtpUseCase) Execute(ctx context.Context, userID uuid.UUID, mobile string) (error) {
	// generate random N-digit OTP

	otp, err := utils.GenerateOTP(uc.strength)
	if err != nil {
		return err
	}

	// OTP valid for N minutes
	expiresAt := utils.NowUTC().Add(5 * time.Minute)

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