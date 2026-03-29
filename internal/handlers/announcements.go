package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type AnnouncementHandler struct {
	service domain.AnnouncementService
}

func NewAnnouncementHandler(service domain.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: service}
}

func (h *AnnouncementHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
		return
	}

	var a domain.Announcement
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid input"}, nil)
		return
	}

	a.CommunityID = communityID
	a.AuthorID = user.ID

	if err := h.service.Create(r.Context(), &a); err != nil {
		h.handleError(w, err)
		return
	}

	domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: a}, nil)
}

func (h *AnnouncementHandler) List(w http.ResponseWriter, r *http.Request) {
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
		return
	}

	announcements, err := h.service.ListActive(r.Context(), communityID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: announcements}, nil)
}

func (h *AnnouncementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
		return
	}

	if err := h.service.Delete(r.Context(), id, communityID); err != nil {
		h.handleError(w, err)
		return
	}
	domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "Deleted"}, nil)
}

func (h *AnnouncementHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*domain.AppError); ok {
		var errList []string
		for k, v := range appErr.Errors {
			errList = append(errList, fmt.Sprintf("%s: %s", k, v))
		}
		domain.WriteJSON(w, appErr.Code, domain.Envelope{Success: false, Message: appErr.Message, Errors: errList}, nil)
		return
	}
	domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Internal server error"}, nil)
}
