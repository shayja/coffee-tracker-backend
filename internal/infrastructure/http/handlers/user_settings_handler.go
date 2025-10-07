// file: internal/infrastructure/http/handlers/user_settings_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/entities"
	http_utils "coffee-tracker-backend/internal/infrastructure/http"
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
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

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
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	var body struct {
		Key   int         `json:"key"`   // enum number
		Value interface{} `json:"value"` // could be bool, string, int
	}
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	setting := entities.Setting(body.Key)
	if !setting.IsValid() {
		http.Error(w, "Invalid setting key", http.StatusBadRequest)
		return
	}

	if err := h.updateUC.Execute(r.Context(), userID, setting, body.Value); err != nil {
		http.Error(w, "Failed to update setting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}