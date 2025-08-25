package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/auth"
	"coffee-tracker-backend/internal/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type AuthHandler struct {
	secret      string
	authService *services.AuthService
    userService *services.UserService
}

type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    User         UserDTO   `json:"user"`
}

type UserDTO struct {
    ID     uuid.UUID `json:"id"`
    Name   string `json:"name"`
    Mobile string `json:"mobile"`
}


func NewAuthHandler(secret string, authService *services.AuthService, userService *services.UserService) *AuthHandler {
	if secret == "" {
		log.Fatal("JWT secret is required")
	}
	return &AuthHandler{
		secret:      secret,
		authService: authService,
        userService: userService,
	}
}

func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Mobile string `json:"mobile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Mobile == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Lookup user by email (ensure they exist)
	user, err := h.userService.GetByMobile(r.Context(), body.Mobile)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	otp, err := h.authService.GenerateOTP(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated OTP for user %s: %s", user.ID, otp)

	// For development, return OTP in response (REMOVE in production!)
	//json.NewEncoder(w).Encode(map[string]string{"otp": otp})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully"})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Mobile string `json:"mobile"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByMobile(r.Context(), body.Mobile)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	valid, err := h.authService.ValidateOTP(r.Context(), user.ID, body.OTP)
	if err != nil || !valid {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Issue JWT + refresh token
	accessToken, err := auth.GenerateJWT(h.secret, user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.authService.GenerateRefreshToken(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}
	//log.Print(accessToken)
	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserDTO{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
		},
	}
	json.NewEncoder(w).Encode(response)
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

	// log.Printf("Refresh token: %s", body.RefreshToken)

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