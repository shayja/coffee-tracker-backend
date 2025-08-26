package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"coffee-tracker-backend/internal/domain/repositories"
	"coffee-tracker-backend/internal/infrastructure/config"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepo repositories.AuthRepository
	cfg  *config.Config
	
}

func NewAuthService(authRepo repositories.AuthRepository, config  *config.Config) *AuthService {
	return &AuthService{authRepo: authRepo, cfg: config}
}

func (s *AuthService) GenerateOTP(ctx context.Context, userID uuid.UUID) (string, error) {
	otp := fmt.Sprintf("%06d", randInt(100000, 999999)) // 6-digit code
	expiresAt := time.Now().Add(5 * time.Minute)

	err := s.authRepo.SaveOTP(ctx, userID.String(), otp, expiresAt)
	if err != nil {
		return "", err
	}

	// TODO: send OTP via email/SMS (integration needed)
	return otp, nil
}

func (s *AuthService) ValidateOTP(ctx context.Context, userID uuid.UUID, otp string) (bool, error) {
	
	// temp: bypass OTP validation in dev mode
	if otp == s.cfg.MagicOtp {
		return true, nil
	}
	// Check if OTP is valid and not expired
	valid, err := s.authRepo.GetValidOTP(ctx, userID.String(), otp)
	if err != nil || !valid {
		return false, err
	}

	// Invalidate OTP after use
	_ = s.authRepo.InvalidateOTP(ctx, userID.String(), otp)

	return true, nil
}


func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
    tokenBytes := make([]byte, 32)
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }

    refreshToken := hex.EncodeToString(tokenBytes)

    expiresAt := time.Now().Add(30 * 24 * time.Hour) // int64
    err := s.authRepo.SaveRefreshToken(ctx, userID, refreshToken, expiresAt)
    if err != nil {
        return "", err
    }

    return refreshToken, nil
}


func (s *AuthService) ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
    userID, err := s.authRepo.GetUserIDByRefreshToken(ctx, token)
    if err != nil {
        if errors.Is(err, repositories.ErrNotFound) {
            return uuid.Nil, errors.New("invalid or expired refresh token")
        }
        return uuid.Nil, err
    }
    return userID, nil
}

func (s *AuthService) RotateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	newToken, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err := s.authRepo.SaveRefreshToken(ctx, userID, newToken, expiresAt); err != nil {
		return "", err
	}
	return newToken, nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// helper for random int
func randInt(min, max int) int {
	b := make([]byte, 1)
	rand.Read(b)
	return min + int(b[0])%(max-min+1)
}

func (s *AuthService) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error{
	return s.authRepo.SaveRefreshToken(ctx, userID, token, expiresAt)
}
func (s *AuthService) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	return s.authRepo.GetRefreshToken(ctx, userID)
}
func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return s.authRepo.DeleteRefreshToken(ctx, userID)
}
