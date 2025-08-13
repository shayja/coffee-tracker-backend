package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"coffee-tracker-backend/internal/usecases"

	"github.com/google/uuid"
)

type CoffeeEntryHandler struct {
	createUseCase    *usecases.CreateCoffeeEntryUseCase
	getEntriesUseCase *usecases.GetCoffeeEntriesUseCase
	getStatsUseCase   *usecases.GetCoffeeStatsUseCase
}

func NewCoffeeEntryHandler(
	createUseCase *usecases.CreateCoffeeEntryUseCase,
	getEntriesUseCase *usecases.GetCoffeeEntriesUseCase,
	getStatsUseCase *usecases.GetCoffeeStatsUseCase,
) *CoffeeEntryHandler {
	return &CoffeeEntryHandler{
		createUseCase:     createUseCase,
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
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req.UserID = userID

	entry, err := h.createUseCase.Execute(r.Context(), req)
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
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	entries, err := h.getEntriesUseCase.Execute(r.Context(), userID, &dateStr, limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *CoffeeEntryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
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
