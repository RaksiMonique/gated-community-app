package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResidentHandler struct {
	service domain.ResidentService
	db      *pgxpool.Pool // Keep for legacy or direct updates not yet refactored
}

func NewResidentHandler(service domain.ResidentService, db *pgxpool.Pool) *ResidentHandler {
	return &ResidentHandler{service: service, db: db}
}

func (h *ResidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
		return
	}

	var res domain.Resident
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "invalid input"}, nil)
		return
	}
	res.CommunityID = communityID

	if err := h.service.Onboard(r.Context(), &res); err != nil {
		h.handleError(w, err)
		return
	}

	domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: res}, nil)
}

func (h *ResidentHandler) List(w http.ResponseWriter, r *http.Request) {
	communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
	if !ok {
		domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "missing community context"}, nil)
		return
	}

	residents, err := h.service.List(r.Context(), communityID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: residents}, nil)
}

func HandleGetResident(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		query := `SELECT id, community_id, unit_id, first_name, last_name, email, phone, resident_type, status, move_in_date, move_out_date FROM residents WHERE id = $1`

		var res domain.Resident
		err := db.QueryRow(r.Context(), query, id).Scan(
			&res.ID, &res.CommunityID, &res.UnitID, &res.FirstName, &res.LastName,
			&res.Email, &res.Phone, &res.ResidentType, &res.Status, &res.MoveInDate, &res.MoveOutDate)

		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "resident not found"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: res}, nil)
	}
}

func HandleSetPrimaryContact(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		tx, err := db.Begin(r.Context())
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "transaction error"}, nil)
			return
		}
		defer tx.Rollback(r.Context())

		// Get household ID for this resident
		var householdID string
		err = tx.QueryRow(r.Context(), "SELECT household_id FROM residents WHERE id = $1", id).Scan(&householdID)
		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "resident not found"}, nil)
			return
		}

		// Reset current primary
		_, err = tx.Exec(r.Context(), "UPDATE residents SET is_primary_contact = false WHERE household_id = $1", householdID)
		if err != nil {
			return
		}

		// Set new primary
		_, err = tx.Exec(r.Context(), "UPDATE residents SET is_primary_contact = true WHERE id = $1", id)
		if err != nil {
			return
		}

		tx.Commit(r.Context())
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "primary contact updated"}, nil)
	}
}

func HandleUpdateResident(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var res domain.Resident
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "invalid input"}, nil)
			return
		}

		query := `
			UPDATE residents 
			SET first_name = $1, last_name = $2, email = $3, phone = $4, resident_type = $5, status = $6, move_in_date = $7, move_out_date = $8
			WHERE id = $9`

		_, err := db.Exec(r.Context(), query,
			res.FirstName, res.LastName, res.Email, res.Phone,
			res.ResidentType, res.Status, res.MoveInDate, res.MoveOutDate, id)

		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to update resident"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "resident updated"}, nil)
	}
}

func HandleDeactivateResident(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		query := `UPDATE residents SET status = 'INACTIVE', move_out_date = NOW() WHERE id = $1`

		_, err := db.Exec(r.Context(), query, id)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to deactivate resident"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "resident deactivated"}, nil)
	}
}

func (h *ResidentHandler) handleError(w http.ResponseWriter, err error) {
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
