// file: internal/infrastructure/http/utils.go
package http

import (
	"bytes"
	"coffee-tracker-backend/internal/contextkeys"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var isDev = os.Getenv("APP_ENV") == "dev" || os.Getenv("APP_ENV") == "development"

// GetPathParam extracts a string path parameter (e.g., {id}) from a request.
// Returns an empty string if not found.
func GetPathParam(r *http.Request, key string) string {
	vars := mux.Vars(r)
	if val, ok := vars[key]; ok {
		return val
	}
	return ""
}

func GetEntryIDByRouteOrAbort(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
	entryIDStr := GetPathParam(r, "id") // get the {id} value
	if entryIDStr == "" {
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


// WriteJSON writes a JSON response and optionally logs it if in dev mode.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if isDev {
		logResponse(status, data)
	}

	json.NewEncoder(w).Encode(data)
}

// WriteError returns a structured JSON error response, and logs it in dev.
func WriteError(w http.ResponseWriter, status int, message string, details ...interface{}) {
	resp := map[string]interface{}{
		"error":   message,
		"status":  status,
		"success": false,
	}

	// Optional detailed debug info in dev mode
	if isDev && len(details) > 0 {
		resp["details"] = details
	}

	WriteJSON(w, status, resp)
}

// LogRequest reads and logs the request body safely (for dev only).
func LogRequest(r *http.Request) {
	if !isDev {
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body for later use

	fmt.Printf(
		"\n[HTTP REQUEST] %s %s\nHeaders: %v\nBody: %s\n\n",
		r.Method, r.URL.Path, r.Header, string(bodyBytes),
	)
}

// logResponse prints the response body (dev only).
func logResponse(status int, data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf(
		"[HTTP RESPONSE] %d at %s\nBody: %s\n\n",
		status, time.Now().Format(time.RFC3339), string(jsonData),
	)
}

func GetUserIpAddress(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	return clientIP
}