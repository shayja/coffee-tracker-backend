// file: internal/infrastructure/http/handlers/tapering_journey_handler.go
package handlers

import (
	"coffee-tracker-backend/internal/contextkeys"
	"coffee-tracker-backend/internal/infrastructure/http/dto"
	"coffee-tracker-backend/internal/infrastructure/utils"
	"coffee-tracker-backend/internal/usecases"
	"encoding/json"
	"net/http"
	"time"
)

type TaperingJourneyHandler struct {
	CreateUC *usecases.CreateTaperingJourneyUseCase
	GetUC    *usecases.GetTaperingJourneysUseCase
	UpdateUC *usecases.UpdateTaperingJourneyUseCase
	DeleteUC *usecases.DeleteTaperingJourneyUseCase
}

func NewTaperingJourneyHandler(c *usecases.CreateTaperingJourneyUseCase, g *usecases.GetTaperingJourneysUseCase, u *usecases.UpdateTaperingJourneyUseCase, d *usecases.DeleteTaperingJourneyUseCase) *TaperingJourneyHandler {
	return &TaperingJourneyHandler{CreateUC: c, GetUC: g, UpdateUC: u, DeleteUC: d}
}

func (h *TaperingJourneyHandler) CreateJourney(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaperingJourneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	statusID := 1//init with active status
	var startedAt time.Time
	if req.StartedAt != nil {
		startedAt = *req.StartedAt
	} else {
		startedAt = time.Now()
	}
	j, err := h.CreateUC.Execute(r.Context(),userID, req.GoalFrequency, req.StartLimit, req.TargetLimit, req.ReductionStep, req.StepPeriod, startedAt, statusID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(j)
}

func (h *TaperingJourneyHandler) GetJourneys(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	journeys, err := h.GetUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(journeys)
}

func (h *TaperingJourneyHandler) UpdateJourney(w http.ResponseWriter, r *http.Request) {
	// load by id, patch, save
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TaperingJourneyHandler) DeleteJourney(w http.ResponseWriter, r *http.Request) {

	// Extract user ID from context (set by auth middleware)
	userID, ok := contextkeys.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Get entry ID from path parameter
	entryID, er := utils.GetEntryIDByRoute(r, w)
	if er != nil {
		return
	}

	err := h.DeleteUC.Execute(r.Context(), entryID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
