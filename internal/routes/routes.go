package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gated-community-api/internal/handlers"
)

func ConfigureRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Standard Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// Health Check
	r.Get("/health", handlers.HandleHealth)

	// API V1 Group
	r.Route("/api/v1", func(r chi.Router) {
		// Future: r.Use(middleware.TenantID)
		r.Get("/communities", func(w http.ResponseWriter, r *http.Request) {
			// Placeholder for Community Handler
		})
	})

	return r
}