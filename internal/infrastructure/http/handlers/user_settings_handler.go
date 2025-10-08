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
	getAllUC *usecases.GetUserSettingsUseCase
	updateUC *usecases.UpdateUserSettingUseCase
}

func NewUserSettingsHandler(
	getAllUC *usecases.GetUserSettingsUseCase,
	updateUC *usecases.UpdateUserSettingUseCase,
) *UserSettingsHandler {
	return &UserSettingsHandler{
		getAllUC: getAllUC,
		updateUC: updateUC,
	}
}

// GET /users/settings
func (h *UserSettingsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	settings, err := h.getAllUC.Execute(r.Context(), userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to load settings", err.Error())
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, map[string]any{
		"settings": settings,
	})
}

// PATCH /settings/:key
func (h *UserSettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok {
		return
	}

	var body struct {
		Key   int         `json:"key"`   // enum number
		Value interface{} `json:"value"` // could be bool, string, int
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	setting := entities.Setting(body.Key)
	if !setting.IsValid() {
		http_utils.WriteError(w, http.StatusBadRequest, "Invalid setting key")
		return
	}

	if err := h.updateUC.Execute(r.Context(), userID, setting, body.Value); err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to update setting", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
