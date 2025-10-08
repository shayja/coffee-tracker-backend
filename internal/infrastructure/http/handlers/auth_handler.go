// file: internal/infrastructure/http/handlers/auth_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/infrastructure/auth"
	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"coffee-tracker-backend/internal/infrastructure/http/dto"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/usecases"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	jwtService            *auth.JWTService
	getUserByIDUC         *usecases.GetUserByIDUseCase
	getUserByMobileUC     *usecases.GetUserByMobileUseCase
	genereteOtpUC         *usecases.GenerateOtpUseCase
	validateOtpUC         *usecases.ValidateOtpUseCase
	saveRefreshTokenUC    *usecases.SaveRefreshTokenUseCase
	getRefreshTokenUC     *usecases.GetRefreshTokenUseCase
	deleteRefreshTokenUC  *usecases.DeleteRefreshTokenUseCase
}

func NewAuthHandler(
	jwtService *auth.JWTService,
	getUserByIDUC *usecases.GetUserByIDUseCase,
	getUserByMobileUC *usecases.GetUserByMobileUseCase,
	genereteOtpUC *usecases.GenerateOtpUseCase,
	validateOtpUC *usecases.ValidateOtpUseCase,
	saveRefreshTokenUC *usecases.SaveRefreshTokenUseCase,
	getRefreshTokenUC *usecases.GetRefreshTokenUseCase,
	deleteRefreshTokenUC *usecases.DeleteRefreshTokenUseCase,
) *AuthHandler {
	if jwtService == nil {
		log.Fatal("JWT service is required")
	}
	return &AuthHandler{
		jwtService:           jwtService,
		getUserByIDUC:        getUserByIDUC,
		getUserByMobileUC:    getUserByMobileUC,
		genereteOtpUC:        genereteOtpUC,
		validateOtpUC:        validateOtpUC,
		saveRefreshTokenUC:   saveRefreshTokenUC,
		getRefreshTokenUC:    getRefreshTokenUC,
		deleteRefreshTokenUC: deleteRefreshTokenUC,
	}
}

// POST /auth/request-otp
func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	http_utils.LogRequest(r)

	var req dto.SendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Mobile == "" {
		http_utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.getUserByMobileUC.Execute(r.Context(), req.Mobile)
	if err != nil {
		http_utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	err = h.genereteOtpUC.Execute(r.Context(), user.ID, req.Mobile)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to generate OTP")
		return
	}
	http_utils.WriteJSON(w, http.StatusOK, dto.SendOtpResponse{
		Message: "OTP sent successfully",
	})
}

// POST /auth/verify-otp
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	http_utils.LogRequest(r)

	var req dto.VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.getUserByMobileUC.Execute(r.Context(), req.Mobile)
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	valid, err := h.validateOtpUC.Execute(r.Context(), user.ID, req.OTP)
	if err != nil || !valid {
		http_utils.WriteError(w, http.StatusUnauthorized, "Invalid or expired OTP")
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	refreshExpiry := time.Now().Add(h.jwtService.RefreshExpiry())
	if err := h.saveRefreshTokenUC.Execute(r.Context(), user.ID, req.DeviceID, refreshToken, refreshExpiry); err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		TokenPair: dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: dto.LoggedInUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	})
}

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	http_utils.LogRequest(r)

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DeviceID == uuid.Nil {
		http_utils.WriteError(w, http.StatusBadRequest, "Missing arguments")
		return
	}

	tokenString, err := utils.ExtractBearerToken(r)
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, "No refresh token found")
		return
	}

	userID, err := h.jwtService.ExtractUserIDFromToken(tokenString)
	if err != nil {
		http_utils.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	claims, err := h.jwtService.ValidateTokenString(tokenString)
	if err != nil || !h.jwtService.IsRefreshToken(claims) {
		http_utils.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	storedToken, expiresAt, err := h.getRefreshTokenUC.Execute(r.Context(), userID, req.DeviceID)
	if err != nil || storedToken != tokenString {
		http_utils.WriteError(w, http.StatusUnauthorized, "Refresh token mismatch")
		return
	}

	if time.Now().After(expiresAt) {
		http_utils.WriteError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	newAccessToken, err := h.jwtService.GenerateAccessToken(userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	newRefreshToken, err := h.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	newRefreshExpiry := time.Now().Add(h.jwtService.RefreshExpiry())
	if err := h.saveRefreshTokenUC.Execute(r.Context(), userID, req.DeviceID, newRefreshToken, newRefreshExpiry); err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, dto.RefreshTokenResponse{
		TokenPair: dto.TokenPair{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
	})
}

// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http_utils.LogRequest(r)

	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	var req dto.DeleteTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DeviceID == uuid.Nil {
		http_utils.WriteError(w, http.StatusBadRequest, "Invalid or missing device_id")
		return
	}

	if err := h.deleteRefreshTokenUC.Execute(r.Context(), userID, req.DeviceID); err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GET /auth/profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	http_utils.LogRequest(r)

	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	user, err := h.getUserByIDUC.Execute(r.Context(), userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, dto.LoggedInUserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Mobile: user.Mobile,
	})
}
