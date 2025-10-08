// file: internal/infrastructure/http/handlers/health_handler.go
package handlers

import (
	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "coffee-tracker-api",
	}

	http_utils.WriteJSON(w, http.StatusOK, response)
}
