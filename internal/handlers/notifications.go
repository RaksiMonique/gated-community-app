package handlers

import (
	"context"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	CommunityID    string    `json:"community_id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	IsRead         bool      `json:"is_read"`
	LinkToResource string    `json:"link_to_resource"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateNotification is an internal helper to trigger notifications from other modules
func CreateNotification(ctx context.Context, db *pgxpool.Pool, userID, communityID, title, message, link string) error {
	query := `
		INSERT INTO notifications (user_id, community_id, title, message, link_to_resource)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(ctx, query, userID, communityID, title, message, link)
	return err
}

func HandleListNotifications(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		query := `
			SELECT id, title, message, is_read, link_to_resource, created_at 
			FROM notifications 
			WHERE user_id = $1 AND community_id = $2
			ORDER BY created_at DESC LIMIT 50`

		rows, err := db.Query(r.Context(), query, user.ID, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch notifications"}, nil)
			return
		}
		defer rows.Close()

		notifications := []Notification{}
		for rows.Next() {
			var n Notification
			rows.Scan(&n.ID, &n.Title, &n.Message, &n.IsRead, &n.LinkToResource, &n.CreatedAt)
			notifications = append(notifications, n)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: notifications}, nil)
	}
}

func HandleMarkNotificationRead(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)

		query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
		_, err := db.Exec(r.Context(), query, id, user.ID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to update notification"}, nil)
			return
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true}, nil)
	}
}

func HandleMarkAllNotificationsRead(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND community_id = $2`
		_, err := db.Exec(r.Context(), query, user.ID, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to update notifications"}, nil)
			return
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true}, nil)
	}
}

func HandleGetUnreadCount(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		var count int
		query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND community_id = $2 AND is_read = false`
		db.QueryRow(r.Context(), query, user.ID, communityID).Scan(&count)

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: count}, nil)
	}
}
