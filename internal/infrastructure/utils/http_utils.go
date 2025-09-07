package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
