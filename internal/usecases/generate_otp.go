// file: internal/usecases/generate_otp.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GenerateOtpUseCase struct {
	authRepo repositories.AuthRepository
}

func NewGenerateOtpUseCase(authRepo repositories.AuthRepository) *GenerateOtpUseCase {
	return &GenerateOtpUseCase{authRepo: authRepo}
}

func (uc *GenerateOtpUseCase) Execute(ctx context.Context, userID uuid.UUID) (string, error) {
	otp := fmt.Sprintf("%06d", randInt(100000, 999999)) // 6-digit code
	expiresAt := time.Now().Add(5 * time.Minute)

	err := uc.authRepo.SaveOTP(ctx, userID, otp, expiresAt)
	if err != nil {
		return "", err
	}

	// TODO: send OTP via email/SMS (integration needed)
	return otp, nil
}

// helper for random int
func randInt(min, max int) int {
	b := make([]byte, 1)
	rand.Read(b)
	return min + int(b[0])%(max-min+1)
}