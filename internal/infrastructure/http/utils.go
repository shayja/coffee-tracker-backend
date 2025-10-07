// file: internal/infrastructure/utils/http_utils.go
package utils

import (
	"coffee-tracker-backend/internal/contextkeys"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// WriteError writes a JSON-formatted error response with the given status and message.
// Example:
//   http_utils.WriteError(w, http.StatusBadRequest, "Invalid request")
func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// WriteJSON writes any Go value as a JSON response with status 200 (OK) by default.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func GetEntryIDByRoute(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
	vars := mux.Vars(r)          // extract path variables
	entryIDStr, ok := vars["id"] // get the {id} value
	if !ok || entryIDStr == "" {
		http.Error(w, "Missing entry ID", http.StatusBadRequest)
		return uuid.UUID{}, nil
	}

	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		http.Error(w, "Invalid entry ID format", http.StatusBadRequest)
		return uuid.UUID{}, nil
	}
	return entryID, err
}
func GetUserIDOrAbort(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
    userID, ok := contextkeys.UserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
    }
    return userID, ok
}