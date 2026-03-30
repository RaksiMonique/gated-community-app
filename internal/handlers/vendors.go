package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Vendor struct {
	ID          string    `json:"id"`
	CommunityID string    `json:"community_id"`
	CompanyName string    `json:"company_name"`
	Category    string    `json:"category"`
	ContactName string    `json:"contact_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Status      string    `json:"status"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type VendorDetails struct {
	Vendor
	JobHistory []MaintenanceRequest `json:"job_history"`
	Stats      struct {
		TotalJobs     int `json:"total_jobs"`
		CompletedJobs int `json:"completed_jobs"`
	} `json:"stats"`
}

func HandleCreateVendor(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		var v Vendor
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid input"}, nil)
			return
		}

		query := `
			INSERT INTO vendors (community_id, company_name, category, contact_name, email, phone, address, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'ACTIVE')
			RETURNING id, created_at`

		err := db.QueryRow(r.Context(), query,
			communityID, v.CompanyName, v.Category, v.ContactName, v.Email, v.Phone, v.Address,
		).Scan(&v.ID, &v.CreatedAt)

		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to create vendor"}, nil)
			return
		}
		domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: v}, nil)
	}
}

func HandleListVendors(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		category := r.URL.Query().Get("category")

		query := `SELECT id, company_name, category, contact_name, email, phone, status, rating FROM vendors WHERE community_id = $1`
		args := []any{communityID}

		if category != "" {
			query += " AND category = $2"
			args = append(args, category)
		}

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Fetch failed"}, nil)
			return
		}
		defer rows.Close()

		vendors := []Vendor{}
		for rows.Next() {
			var v Vendor
			rows.Scan(&v.ID, &v.CompanyName, &v.Category, &v.ContactName, &v.Email, &v.Phone, &v.Status, &v.Rating)
			vendors = append(vendors, v)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: vendors}, nil)
	}
}

func HandleGetVendor(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var details VendorDetails

		// Get Vendor Info
		err := db.QueryRow(r.Context(),
			"SELECT id, company_name, category, contact_name, email, phone, address, status, rating FROM vendors WHERE id = $1 AND community_id = $2",
			id, communityID).Scan(&details.ID, &details.CompanyName, &details.Category, &details.ContactName, &details.Email, &details.Phone, &details.Address, &details.Status, &details.Rating)

		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "Vendor not found"}, nil)
			return
		}

		// Get Job History & Stats
		rows, _ := db.Query(r.Context(),
			"SELECT id, title, status, priority, created_at FROM maintenance_requests WHERE vendor_id = $1 ORDER BY created_at DESC LIMIT 10", id)
		defer rows.Close()

		for rows.Next() {
			var mr MaintenanceRequest
			rows.Scan(&mr.ID, &mr.Title, &mr.Status, &mr.Priority, &mr.CreatedAt)
			details.JobHistory = append(details.JobHistory, mr)
			details.Stats.TotalJobs++
			if mr.Status == "COMPLETED" {
				details.Stats.CompletedJobs++
			}
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: details}, nil)
	}
}

func HandleDeleteVendor(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		_, err := db.Exec(r.Context(), "DELETE FROM vendors WHERE id = $1 AND community_id = $2", id, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Delete failed"}, nil)
			return
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Message: "Vendor removed"}, nil)
	}
}
