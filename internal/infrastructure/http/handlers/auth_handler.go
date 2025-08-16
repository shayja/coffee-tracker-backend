package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

type AuthHandler struct {
	secret      string
	authService *services.AuthService
}

func NewAuthHandler(secret string, authService *services.AuthService) *AuthHandler {
	if secret == "" {
		log.Fatal("JWT secret is required")
	}
	return &AuthHandler{
		secret:      secret,
		authService: authService,
	}
}

func (h *AuthHandler) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context
    userID, ok := contextkeys.UserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Generate access token (short-lived JWT)
    accessToken, err := auth.GenerateJWT(h.secret, userID)
    if err != nil {
        log.Printf("failed to generate JWT: %v", err)
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Generate refresh token (long-lived, persisted in DB)
    refreshToken, err := h.authService.GenerateRefreshToken(r.Context(), userID)
    if err != nil {
        log.Printf("failed to generate refresh token: %v", err)
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Return both tokens
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }); err != nil {
        log.Printf("failed to encode response: %v", err)
    }
}

// Issue a new access token for authenticated user
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    type requestBody struct {
        RefreshToken string `json:"refresh_token"`
    }

    var body requestBody
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
        http.Error(w, "Missing refresh token", http.StatusUnauthorized)
        return
    }

    // Validate & get user ID
    userID, err := h.authService.ValidateRefreshToken(r.Context(), body.RefreshToken)
    if err != nil {
        http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
        return
    }

    // Issue new tokens
    newAccessToken, err := auth.GenerateJWT(h.secret, userID)
    if err != nil {
        http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
        return
    }

    newRefreshToken, err := h.authService.RotateRefreshToken(r.Context(), userID)
    if err != nil {
        http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token":  newAccessToken,
        "refresh_token": newRefreshToken,
    })
}