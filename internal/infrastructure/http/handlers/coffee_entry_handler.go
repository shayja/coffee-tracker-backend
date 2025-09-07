// file: internal/infrastructure/http/handlers/coffee_entry_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/usecases"
)

type CoffeeEntryHandler struct {
	createCoffeeUC    *usecases.CreateCoffeeEntryUseCase
	editCoffeeUC      *usecases.EditCoffeeEntryUseCase
	deleteCoffeeUC     *usecases.DeleteCoffeeEntryUseCase
	listCoffeeUC *usecases.ListCoffeeEntriesUseCase
	getStatsUseCase   *usecases.GetCoffeeStatsUseCase
}

func NewCoffeeEntryHandler(
	createCoffeeUC *usecases.CreateCoffeeEntryUseCase,
	editCoffeeUC *usecases.EditCoffeeEntryUseCase,
	deleteCoffeeUC *usecases.DeleteCoffeeEntryUseCase,
	listCoffeeUC *usecases.ListCoffeeEntriesUseCase,
	getStatsUseCase *usecases.GetCoffeeStatsUseCase,
) *CoffeeEntryHandler {
	return &CoffeeEntryHandler{
		createCoffeeUC:     createCoffeeUC,
		editCoffeeUC:       editCoffeeUC,
		deleteCoffeeUC:     deleteCoffeeUC,
		listCoffeeUC: listCoffeeUC,
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

	entry, err := h.createCoffeeUC.Execute(r.Context(), req, userID)
	if err != nil {
		switch err {
		case usecases.ErrInvalidInput:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecases.ErrConflict:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			// Log the error for debugging
			log.Printf("CreateEntry error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Set Location header with URL to the newly created resource
	location := fmt.Sprintf("%s/%s", strings.TrimSuffix(r.URL.Path, "/"), entry.ID.String())
	w.Header().Set("Location", location)
	
	// Return 201 Created status
	w.WriteHeader(http.StatusCreated)
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
	entryID, err := utils.GetEntryIDByRoute(r, w)
	if err != nil {
		return
	}
	req.ID = entryID

	entry, err := h.editCoffeeUC.Execute(r.Context(), req, userID)
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

	entries, err := h.listCoffeeUC.Execute(r.Context(), userID, &dateStr, tzOffset, limit, offset)
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
	entryID, err := utils.GetEntryIDByRoute(r, w)
	if err != nil {
		return
	}

	// Call use case to delete
	err = h.deleteCoffeeUC.Execute(r.Context(), userID, entryID)
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

