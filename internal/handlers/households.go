package handlers

import (
	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"
	"net/http"
)

type HouseholdHandler struct {
	service domain.HouseholdService
}

func NewHouseholdHandler(service domain.HouseholdService) *HouseholdHandler {
	return &HouseholdHandler{service: service}
}

func (h *HouseholdHandler) List(w http.ResponseWriter, r *http.Request) {
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "missing community context"}, nil)
		return
	}

	households, err := h.service.List(r.Context(), communityID)
	if err != nil {
		domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to fetch households"}, nil)
		return
	}

	domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: households}, nil)
}

func (h *HouseholdHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Implementation for manual household creation
	domain.WriteJSON(w, http.StatusNotImplemented, domain.Envelope{Success: false, Message: "Manual creation not yet implemented"}, nil)
}
