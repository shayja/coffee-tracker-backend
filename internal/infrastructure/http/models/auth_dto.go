// file: internal/infrastructure/http/models/auth_dto.go
package models

import "github.com/google/uuid"

// SendOtpRequest represents a request to send an OTP to the user's mobile number.
type SendOtpRequest struct {
    Mobile string `json:"mobile" binding:"required"`
}

// VerifyOtpRequest is used to verify a previously sent OTP.
type VerifyOtpRequest struct {
    Mobile   string    `json:"mobile" binding:"required"`
    OTP      string    `json:"otp" binding:"required,min=6,max=6"`
    DeviceID uuid.UUID `json:"device_id" binding:"required"`
}

type DeleteTokenRequest struct {
	DeviceID	uuid.UUID  `json:"device_id" binding:"required"`
}

type RefreshTokenRequest struct {
	DeviceID	uuid.UUID  `json:"device_id" binding:"required"`
}

// Response DTOs
type SendOtpResponse struct {
	Message  string `json:"message"`
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
    TokenPair
    User LoggedInUserResponse `json:"user"`
}

type RefreshTokenResponse struct {
    TokenPair
}