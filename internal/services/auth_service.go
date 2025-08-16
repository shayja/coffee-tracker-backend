package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"coffee-tracker-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type AuthService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
    tokenBytes := make([]byte, 32)
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }

    refreshToken := hex.EncodeToString(tokenBytes)

    expiresAt := time.Now().Add(30 * 24 * time.Hour).Unix() // int64
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
	if err := s.authRepo.SaveRefreshToken(ctx, userID, newToken, expiresAt.Unix()); err != nil {
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
