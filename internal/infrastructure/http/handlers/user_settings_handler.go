// file: internal/infrastructure/http/handlers/user_settings_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/usecases"
	"encoding/json"
	"net/http"
)

type UserSettingsHandler struct {
	getAllUC   *usecases.GetUserSettingsUseCase
	updateUC   *usecases.UpdateUserSettingUseCase
}

func NewUserSettingsHandler(getAllUC *usecases.GetUserSettingsUseCase, updateUC *usecases.UpdateUserSettingUseCase) *UserSettingsHandler {
	return &UserSettingsHandler{
		getAllUC: getAllUC,
		updateUC: updateUC,
	}
}

// GET /users/settings
func (h *UserSettingsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	settings, err := h.getAllUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to load settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"settings": settings,
	})
}

// PUT /users/settings/:key
func (h *UserSettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	key := r.PathValue("key") // if you’re using Go 1.22+ net/http patterns
	if key == "" {
		http.Error(w, "Missing setting key", http.StatusBadRequest)
		return
	}

	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := usecases.UpdateUserSettingRequest{Key: key, Value: body.Value}

	if err := h.updateUC.Execute(r.Context(), userID, req); err != nil {
		http.Error(w, "Failed to update setting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
