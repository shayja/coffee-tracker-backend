// file: internal/infrastructure/http/handlers/generic_kv_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"coffee-tracker-backend/internal/usecases"
)

type GenericKVHandler struct {
    useCase *usecases.GetGenericKVUseCase
}

func NewGenericKVHandler(useCase *usecases.GetGenericKVUseCase) *GenericKVHandler {
    return &GenericKVHandler{useCase: useCase}
}

func (h *GenericKVHandler) Get(w http.ResponseWriter, r *http.Request) {
    // Read query parameters
    typeStr := r.URL.Query().Get("type")
    language := r.URL.Query().Get("language")

    typeID, err := strconv.Atoi(typeStr)
    if err != nil || typeStr == "" {
        http.Error(w, "invalid or missing 'type' parameter", http.StatusBadRequest)
        return
    }
    if language == "" {
        http.Error(w, "missing 'language' parameter", http.StatusBadRequest)
        return
    }

    items, err := h.useCase.Execute(r.Context(), typeID, language)
    if err != nil {
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "items": items,
    })
}
