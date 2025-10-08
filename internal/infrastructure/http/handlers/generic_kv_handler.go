// file: internal/infrastructure/http/handlers/generic_kv_handler.go
package handlers

import (
	"net/http"
	"strconv"

	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"coffee-tracker-backend/internal/usecases"
)

type GenericKVHandler struct {
	useCase *usecases.GetGenericKVUseCase
}

func NewGenericKVHandler(useCase *usecases.GetGenericKVUseCase) *GenericKVHandler {
	return &GenericKVHandler{useCase: useCase}
}

// GET /kv?type=123&language=en
func (h *GenericKVHandler) Get(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")
	language := r.URL.Query().Get("language")

	if typeStr == "" {
		http_utils.WriteError(w, http.StatusBadRequest, "missing 'type' parameter")
		return
	}
	typeID, err := strconv.Atoi(typeStr)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "'type' must be a number")
		return
	}

	if language == "" {
		http_utils.WriteError(w, http.StatusBadRequest, "missing 'language' parameter")
		return
	}

	items, err := h.useCase.Execute(r.Context(), typeID, language)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to get items", err.Error())
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
