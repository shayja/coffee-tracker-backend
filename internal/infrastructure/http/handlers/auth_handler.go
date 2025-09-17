// file: internal/infrastructure/http/handlers/auth_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/infrastructure/http/dto"
	"coffee-tracker-backend/internal/usecases"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	jwtService  *auth.JWTService
	getUserByIDUC *usecases.GetUserByIDUseCase
	getUserByMobileUC *usecases.GetUserByMobileUseCase
	genereteOtpUC *usecases.GenerateOtpUseCase
	validateOtpUC *usecases.ValidateOtpUseCase
	saveRefreshTokenUC *usecases.SaveRefreshTokenUseCase
	getRefreshTokenUC *usecases.GetRefreshTokenUseCase
	deleteRefreshTokenUC *usecases.DeleteRefreshTokenUseCase
}

func NewAuthHandler(
	jwtService *auth.JWTService, 
	getUserByIDUC *usecases.GetUserByIDUseCase, 
	getUserByMobileUC *usecases.GetUserByMobileUseCase, 
	genereteOtpUC *usecases.GenerateOtpUseCase,
	validateOtpUC *usecases.ValidateOtpUseCase,
	saveRefreshTokenUC *usecases.SaveRefreshTokenUseCase,
	getRefreshTokenUC *usecases.GetRefreshTokenUseCase,
	deleteRefreshTokenUC *usecases.DeleteRefreshTokenUseCase) *AuthHandler {
	if jwtService == nil {
		log.Fatal("JWT service is required")
	}
	return &AuthHandler{
		jwtService:  jwtService,
		getUserByIDUC: getUserByIDUC,
		getUserByMobileUC: getUserByMobileUC,
		genereteOtpUC: genereteOtpUC,
		validateOtpUC: validateOtpUC,
		saveRefreshTokenUC: saveRefreshTokenUC,
		getRefreshTokenUC: getRefreshTokenUC,
		deleteRefreshTokenUC: deleteRefreshTokenUC,
	}
}


func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Mobile == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Lookup user by mobile (ensure they exist)
	user, err := h.getUserByMobileUC.Execute(r.Context(), req.Mobile)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	otp, err := h.genereteOtpUC.Execute(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated OTP for user %s: %s", user.ID, otp)

	w.WriteHeader(http.StatusOK)

	// Return response
	response := dto.SendOtpResponse{
		Message: "OTP sent successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.getUserByMobileUC.Execute(r.Context(), req.Mobile)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Validate OTP
	valid, err := h.validateOtpUC.Execute(r.Context(), user.ID, req.OTP)
	if err != nil || !valid {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Save refresh token to database
	refreshExpiry := time.Now().Add(h.jwtService.RefreshExpiry())
	if err := h.saveRefreshTokenUC.Execute(r.Context(), user.ID, refreshToken, refreshExpiry); err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.LoggedInUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Check if it's a refresh token
	if !h.jwtService.IsRefreshToken(claims) {
		http.Error(w, "Not a refresh token", http.StatusUnauthorized)
		return
	}

	// Extract user ID
	userID, err := h.jwtService.ExtractUserIDFromToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Verify refresh token exists in database and matches
	storedToken, expiresAt, err := h.getRefreshTokenUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	if storedToken != req.RefreshToken {
		http.Error(w, "Refresh token mismatch", http.StatusUnauthorized)
		return
	}

	if time.Now().After(expiresAt) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	// Generate new tokens
	newAccessToken, err := h.jwtService.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := h.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Save new refresh token (rotate)
	newRefreshExpiry := time.Now().Add(h.jwtService.RefreshExpiry())
	if err := h.saveRefreshTokenUC.Execute(r.Context(), userID, newRefreshToken, newRefreshExpiry); err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	// Return new tokens
	response := dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete refresh token from database
	if err := h.deleteRefreshTokenUC.Execute(r.Context(), userID); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate access token
	accessToken, err := h.jwtService.GenerateAccessToken(userID)
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		log.Printf("failed to generate refresh token: %v", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Save refresh token to database
	refreshExpiry := time.Now().Add(h.jwtService.RefreshExpiry())
	if err := h.saveRefreshTokenUC.Execute(r.Context(), userID, refreshToken, refreshExpiry); err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	// Return both tokens
	response := dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Optional: Get current user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.getUserByIDUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := dto.LoggedInUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Mobile:    user.Mobile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}