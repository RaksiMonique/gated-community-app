package handlers

import (
	"net/http"

	"gated-community-api/internal/domain"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	domain.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
