// file: internal/infrastructure/http/handlers/user_settings_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/usecases"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

// PATCH /settings/:key
func (h *UserSettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "Missing setting key", http.StatusBadRequest)
		return
	}

	setting := entities.Setting(key)
	if !setting.IsValid() {
		http.Error(w, "Invalid setting key", http.StatusBadRequest)
		return
	}

	var body struct {
		Value interface{} `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert userID to uuid.UUID
	uid, err := uuid.Parse(userID.String())
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	if err := h.updateUC.Execute(r.Context(), uid, setting, body.Value); err != nil {
		http.Error(w, "Failed to update setting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}