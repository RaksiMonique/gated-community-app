package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VisitorPass struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	VehicleNumber   string    `json:"vehicle_number"`
	ExpectedArrival time.Time `json:"expected_arrival"`
	Purpose         string    `json:"purpose"`
	InviteCode      string    `json:"invite_code"`
	Status          string    `json:"status"`
	UnitID          string    `json:"unit_id"`
}

// generateSecureToken creates a short unique code for the QR pass
func generateSecureToken(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed ambiguous chars like 0, O, 1, I
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func HandleCreateVisitorPass(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var input VisitorPass
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "invalid input"}, nil)
			return
		}

		input.InviteCode = generateSecureToken(8)

		query := `
			INSERT INTO visitors (community_id, unit_id, name, phone, vehicle_number, expected_arrival, invite_code, purpose, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, status`

		err := db.QueryRow(r.Context(), query,
			communityID, input.UnitID, input.Name, input.Phone, input.VehicleNumber,
			input.ExpectedArrival, input.InviteCode, input.Purpose, user.ID,
		).Scan(&input.ID, &input.Status)

		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to create pass"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: input}, nil)
	}
}

func HandleListUpcomingVisitors(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		query := `
			SELECT id, name, vehicle_number, expected_arrival, status, invite_code, purpose
			FROM visitors
			WHERE community_id = $1 AND created_by = $2 AND status = 'SCHEDULED'
			ORDER BY expected_arrival ASC`

		rows, err := db.Query(r.Context(), query, communityID, user.ID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to fetch visitors"}, nil)
			return
		}
		defer rows.Close()

		visitors := []VisitorPass{}
		for rows.Next() {
			var v VisitorPass
			if err := rows.Scan(&v.ID, &v.Name, &v.VehicleNumber, &v.ExpectedArrival, &v.Status, &v.InviteCode, &v.Purpose); err != nil {
				continue
			}
			visitors = append(visitors, v)
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: visitors}, nil)
	}
}

func HandleGetVisitorPass(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var v VisitorPass
		query := `
			SELECT id, name, phone, vehicle_number, expected_arrival, status, invite_code, purpose, unit_id
			FROM visitors
			WHERE id = $1 AND community_id = $2`

		err := db.QueryRow(r.Context(), query, id, communityID).Scan(
			&v.ID, &v.Name, &v.Phone, &v.VehicleNumber, &v.ExpectedArrival, &v.Status, &v.InviteCode, &v.Purpose, &v.UnitID,
		)
		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "pass not found"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: v}, nil)
	}
}

func HandleValidateVisitorPass(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var input struct {
			InviteCode string `json:"invite_code"`
			ActionType string `json:"action_type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "invalid input"}, nil)
			return
		}

		tx, err := db.Begin(r.Context())
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to start transaction"}, nil)
			return
		}
		defer tx.Rollback(r.Context())

		var visitorID, unitID, status string
		var expectedArrival time.Time

		query := `SELECT id, unit_id, status, expected_arrival FROM visitors WHERE invite_code = $1 AND community_id = $2`
		err = tx.QueryRow(r.Context(), query, input.InviteCode, communityID).Scan(&visitorID, &unitID, &status, &expectedArrival)
		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "Pass not found or invalid code"}, nil)
			return
		}

		newStatus := ""
		if input.ActionType == "ENTRY" {
			if status != "SCHEDULED" {
				domain.WriteJSON(w, http.StatusConflict, domain.Envelope{Success: false, Message: "Pass has already been used or is invalid"}, nil)
				return
			}
			// Expiry check: Rejection if more than 24 hours past arrival window
			if time.Now().After(expectedArrival.Add(24 * time.Hour)) {
				domain.WriteJSON(w, http.StatusGone, domain.Envelope{Success: false, Message: "Pass has expired"}, nil)
				return
			}
			newStatus = "ARRIVED"
		} else if input.ActionType == "EXIT" {
			if status != "ARRIVED" {
				domain.WriteJSON(w, http.StatusConflict, domain.Envelope{Success: false, Message: "Visitor is not currently inside"}, nil)
				return
			}
			newStatus = "DEPARTED"
		} else {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid action type"}, nil)
			return
		}

		// Update visitor status
		_, err = tx.Exec(r.Context(), "UPDATE visitors SET status = $1 WHERE id = $2", newStatus, visitorID)
		if err != nil {
			return
		}

		if input.ActionType == "ENTRY" {
			// Create new log session
			_, err = tx.Exec(r.Context(), `
				INSERT INTO security_logs (community_id, person_type, visitor_id, unit_id, time_in, guard_id)
				VALUES ($1, 'VISITOR', $2, $3, NOW(), $4)`,
				communityID, visitorID, unitID, user.ID,
			)
		} else {
			// Close existing log session
			_, err = tx.Exec(r.Context(), `
				UPDATE security_logs SET time_out = NOW() 
				WHERE visitor_id = $1 AND time_out IS NULL`,
				visitorID,
			)
		}

		if err != nil {
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{
			Success: true,
			Message: fmt.Sprintf("Access Authorized: %s recorded", input.ActionType),
		}, nil)
	}
}
