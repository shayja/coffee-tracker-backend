// file: internal/infrastructure/http/dto/auth_dto.go
package dto

// Request DTOs
type SendOtpRequest struct {
	Mobile string `json:"mobile" binding:"required,e164"`
}

type VerifyOtpRequest struct {
	Mobile string `json:"mobile" binding:"required,e164"`
	OTP    string `json:"otp" binding:"required,min=6,max=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Response DTOs
type SendOtpResponse struct {
	Message  string `json:"message"`
}

type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         LoggedInUserResponse `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}