package handlers

import (
	"gated-community-api/internal/domain"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleHealth(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "disconnected"
		if db != nil {
			if err := db.Ping(r.Context()); err == nil {
				dbStatus = "connected"
			}
		}

		payload := domain.Envelope{
			Success: true,
			Data: map[string]string{
				"status":   "ok",
				"database": dbStatus,
			},
		}

		err := domain.WriteJSON(w, http.StatusOK, payload, nil)
		if err != nil {
			log.Printf("failed to write health response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
