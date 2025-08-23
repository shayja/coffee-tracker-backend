package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CoffeeEntryHandler struct {
	createUseCase    *usecases.CreateCoffeeEntryUseCase
	editUseCase      *usecases.EditCoffeeEntryUseCase
	deleteUseCase    *usecases.DeleteCoffeeEntryUseCase
	getEntriesUseCase *usecases.GetCoffeeEntriesUseCase
	getStatsUseCase   *usecases.GetCoffeeStatsUseCase
}

func NewCoffeeEntryHandler(
	createUseCase *usecases.CreateCoffeeEntryUseCase,
	editUseCase *usecases.EditCoffeeEntryUseCase,
	deleteUseCase *usecases.DeleteCoffeeEntryUseCase,
	getEntriesUseCase *usecases.GetCoffeeEntriesUseCase,
	getStatsUseCase *usecases.GetCoffeeStatsUseCase,
) *CoffeeEntryHandler {
	return &CoffeeEntryHandler{
		createUseCase:     createUseCase,
		editUseCase:       editUseCase,
		deleteUseCase:     deleteUseCase,
		getEntriesUseCase: getEntriesUseCase,
		getStatsUseCase:   getStatsUseCase,
	}
}

func (h *CoffeeEntryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req usecases.CreateCoffeeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	entry, err := h.createUseCase.Execute(r.Context(), req, userID)
	if err != nil {
		switch err {
		case usecases.ErrInvalidInput:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *CoffeeEntryHandler) EditEntry(w http.ResponseWriter, r *http.Request) {
	var req usecases.EditCoffeeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get entry ID from path parameter
	entryID, err := getEntryIDByRoute(r, w)
	if err != nil {
		return
	}
	req.ID = entryID

	entry, err := h.editUseCase.Execute(r.Context(), req, userID)
	if err != nil {
		switch err {
		case usecases.ErrInvalidInput:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *CoffeeEntryHandler) GetEntries(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	tzOffsetStr := r.URL.Query().Get("tzOffset")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)


	var tzOffset *int
	if tzOffsetStr != "" {
		if offset, err := strconv.Atoi(tzOffsetStr); err == nil {
			tzOffset = &offset
		}
	}

	entries, err := h.getEntriesUseCase.Execute(r.Context(), userID, &dateStr, tzOffset, limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *CoffeeEntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get entry ID from path parameter
	entryID, err := getEntryIDByRoute(r, w)
	if err != nil {
		return
	}

	// Call use case to delete
	err = h.deleteUseCase.Execute(r.Context(), userID, entryID)
	if err != nil {
		switch err {
		case usecases.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CoffeeEntryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := h.getStatsUseCase.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func getEntryIDByRoute(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
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
