package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceRequest struct {
	ID               string    `json:"id"`
	CommunityID      string    `json:"community_id"`
	UnitID           string    `json:"unit_id"`
	ResidentID       string    `json:"resident_id"`
	VendorID         *string   `json:"vendor_id"`
	AssignedToUserID *string   `json:"assigned_to_user_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Status           string    `json:"status"`
	Priority         string    `json:"priority"`
	Attachments      []string  `json:"attachments"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func HandleCreateMaintenanceRequest(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var req MaintenanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid input"}, nil)
			return
		}

		// Fetch the resident's primary unit and resident ID automatically
		queryResident := `SELECT id, unit_id FROM residents WHERE user_id = $1 AND community_id = $2 AND status = 'ACTIVE' LIMIT 1`
		err := db.QueryRow(r.Context(), queryResident, user.ID, communityID).Scan(&req.ResidentID, &req.UnitID)
		if err != nil {
			domain.WriteJSON(w, http.StatusForbidden, domain.Envelope{Success: false, Message: "User is not an active resident"}, nil)
			return
		}

		query := `
			INSERT INTO maintenance_requests (community_id, unit_id, resident_id, title, description, category, priority, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'OPEN')
			RETURNING id, created_at`

		err = db.QueryRow(r.Context(), query,
			communityID, req.UnitID, req.ResidentID, req.Title, req.Description, req.Category, req.Priority,
		).Scan(&req.ID, &req.CreatedAt)

		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to create request"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: req}, nil)
	}
}

func HandleListMaintenanceRequests(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		status := r.URL.Query().Get("status")
		priority := r.URL.Query().Get("priority")

		query := `SELECT id, unit_id, resident_id, vendor_id, assigned_to_user_id, title, description, category, status, priority, created_at 
		          FROM maintenance_requests WHERE community_id = $1`
		args := []any{communityID}

		if status != "" {
			args = append(args, status)
			query += fmt.Sprintf(" AND status = $%d", len(args))
		}
		if priority != "" {
			args = append(args, priority)
			query += fmt.Sprintf(" AND priority = $%d", len(args))
		}
		query += " ORDER BY created_at DESC"

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch requests"}, nil)
			return
		}
		defer rows.Close()

		requests := []MaintenanceRequest{}
		for rows.Next() {
			var req MaintenanceRequest
			err := rows.Scan(&req.ID, &req.UnitID, &req.ResidentID, &req.VendorID, &req.AssignedToUserID, &req.Title, &req.Description, &req.Category, &req.Status, &req.Priority, &req.CreatedAt)
			if err != nil {
				continue
			}
			requests = append(requests, req)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: requests}, nil)
	}
}

func HandleUpdateMaintenanceStatus(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var input struct {
			Status string `json:"status"`
		}
		json.NewDecoder(r.Body).Decode(&input)

		query := `UPDATE maintenance_requests SET status = $1, updated_at = NOW() WHERE id = $2 AND community_id = $3`
		_, err := db.Exec(r.Context(), query, input.Status, id, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Update failed"}, nil)
			return
		}

		// Trigger Notification to the resident
		var residentUserID string
		err = db.QueryRow(r.Context(),
			`SELECT r.user_id 
			 FROM residents r 
			 JOIN maintenance_requests mr ON r.id = mr.resident_id 
			 WHERE mr.id = $1 AND r.user_id IS NOT NULL`,
			id).Scan(&residentUserID)

		if err == nil {
			CreateNotification(r.Context(), db, residentUserID, communityID,
				"Maintenance Update", "Your request status has changed to: "+input.Status, "/maintenance/"+id)
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "Status updated"}, nil)
	}
}

func HandleAssignMaintenance(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var input struct {
			VendorID *string `json:"vendor_id"`
			StaffID  *string `json:"staff_id"`
		}
		json.NewDecoder(r.Body).Decode(&input)

		query := `
			UPDATE maintenance_requests 
			SET vendor_id = $1, assigned_to_user_id = $2, status = 'IN_PROGRESS', updated_at = NOW() 
			WHERE id = $3 AND community_id = $4`

		_, err := db.Exec(r.Context(), query, input.VendorID, input.StaffID, id, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Assignment failed"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "Request assigned and moved to In Progress"}, nil)
	}
}
