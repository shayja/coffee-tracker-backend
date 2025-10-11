package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	http_utils "coffee-tracker-backend/internal/infrastructure/http"
	"coffee-tracker-backend/internal/infrastructure/http/models"
	"coffee-tracker-backend/internal/usecases"
)

type CoffeeEntryHandler struct {
	createUC    *usecases.CreateCoffeeEntryUseCase
	getAllUC    *usecases.GetCoffeeEntriesUseCase
	updateUC    *usecases.UpdateCoffeeEntryUseCase
	deleteUC    *usecases.DeleteCoffeeEntryUseCase
	clearUC     *usecases.ClearCoffeeEntriesUseCase
	getStatsUC  *usecases.GetCoffeeStatsUseCase
}

func NewCoffeeEntryHandler(
	createUC *usecases.CreateCoffeeEntryUseCase,
	getAllUC *usecases.GetCoffeeEntriesUseCase,
	updateUC *usecases.UpdateCoffeeEntryUseCase,
	deleteUC *usecases.DeleteCoffeeEntryUseCase,
	clearUC *usecases.ClearCoffeeEntriesUseCase,
	getStatsUC *usecases.GetCoffeeStatsUseCase,
) *CoffeeEntryHandler {
	return &CoffeeEntryHandler{
		createUC:   createUC,
		getAllUC:   getAllUC,
		updateUC:   updateUC,
		deleteUC:   deleteUC,
		clearUC:    clearUC,
		getStatsUC: getStatsUC,
	}
}

// POST /entries
func (h *CoffeeEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	var req models.CreateCoffeeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := h.createUC.Execute(r.Context(), userID, &req)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http_utils.WriteJSON(w, http.StatusCreated, entry)
}

// GET /entries
func (h *CoffeeEntryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

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

	entries, err := h.getAllUC.Execute(r.Context(), userID, &dateStr, tzOffset, limit, offset)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to get entries")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, entries)
}

// PATCH /entries/{id}
func (h *CoffeeEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }


	entryID, err := http_utils.GetEntryIDByRouteOrAbort(r, w)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, usecases.ErrInvalidInput.Error())
		return
	}

	var req models.UpdateCoffeeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := h.updateUC.Execute(r.Context(), userID, entryID, &req)
	if err != nil {
		switch err {
		case usecases.ErrNotFound:
			http_utils.WriteError(w, http.StatusNotFound, err.Error())
		case usecases.ErrInvalidInput:
			http_utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			http_utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, entry)
}

// DELETE /entries/{id}
func (h *CoffeeEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	entryID, err := http_utils.GetEntryIDByRouteOrAbort(r, w)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, usecases.ErrInvalidInput.Error())
		return
	}

	err = h.deleteUC.Execute(r.Context(), userID, entryID)
	if err != nil {
		switch err {
		case usecases.ErrNotFound:
			http_utils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			http_utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /entries
func (h *CoffeeEntryHandler) ClearAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	err := h.clearUC.Execute(r.Context(), userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /entries/stats
func (h *CoffeeEntryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := http_utils.GetUserIDOrAbort(w, r)
	if !ok { return }

	stats, err := h.getStatsUC.Execute(r.Context(), userID)
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	http_utils.WriteJSON(w, http.StatusOK, stats)
}
