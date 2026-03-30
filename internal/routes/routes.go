package routes

import (
	"gated-community-api/internal/handlers"
	"gated-community-api/internal/middleware"
	"gated-community-api/internal/repositories"
	"gated-community-api/internal/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Note: In a real app, jwtSecret should come from config
const jwtSecret = "your-very-secure-secret"

func ConfigureRoutes(db *pgxpool.Pool) *chi.Mux {
	// Dependency Injection Setup
	annRepo := repositories.NewPostgresAnnouncementRepository(db)
	annService := services.NewAnnouncementService(annRepo)
	annHandler := handlers.NewAnnouncementHandler(annService)

	resRepo := repositories.NewPostgresResidentRepository(db)
	resService := services.NewResidentService(resRepo)
	resHandler := handlers.NewResidentHandler(resService, db)

	houseRepo := repositories.NewPostgresHouseholdRepository(db)
	houseService := services.NewHouseholdService(houseRepo)
	houseHandler := handlers.NewHouseholdHandler(houseService)

	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// API Version 1
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handlers.HandleHealth(db))
		r.Post("/auth/login", handlers.HandleLogin(db, jwtSecret))

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtSecret))
			// Multi-tenancy enforcement middleware
			r.Use(middleware.TenantScope(db))

			r.Get("/auth/me", handlers.HandleGetMe(db))

			// Security Logs
			r.With(middleware.RequireRole(db, "ADMIN", "SECURITY")).
				Get("/security/logs", handlers.HandleListSecurityLogs(db))

			// Visitor Pass Management
			r.Route("/visitors", func(r chi.Router) {
				r.Get("/", handlers.HandleListUpcomingVisitors(db))
				r.Post("/", handlers.HandleCreateVisitorPass(db))
				r.Get("/{id}", handlers.HandleGetVisitorPass(db))
			})

			// Security Specific Routes
			r.With(middleware.RequireRole(db, "SECURITY")).
				Post("/visitors/validate", handlers.HandleValidateVisitorPass(db))

			// Example of role-based route
			r.With(middleware.RequireRole(db, "ADMIN")).Get("/admin/dashboard", handlers.HandleGetAdminDashboard(db))
			r.With(middleware.RequireRole(db, "ADMIN")).Get("/admin/stats", handlers.HandleHealth(db))

			// Resident Management
			r.Route("/residents", func(r chi.Router) {
				r.Get("/", resHandler.List)
				r.Post("/", resHandler.Create)
				r.Get("/{id}", handlers.HandleGetResident(db))
				r.Put("/{id}", handlers.HandleUpdateResident(db))
				r.Patch("/{id}/primary", handlers.HandleSetPrimaryContact(db))
				r.Delete("/{id}", handlers.HandleDeactivateResident(db))
			})

			// Household Management
			r.Route("/households", func(r chi.Router) {
				r.Get("/", houseHandler.List)
				r.Post("/", houseHandler.Create)
			})

			// Billing Module
			r.Route("/billing", func(r chi.Router) {
				r.With(middleware.RequireRole(db, "ADMIN")).Post("/invoices", handlers.HandleCreateInvoice(db))
				r.Get("/invoices", handlers.HandleListInvoices(db))
				r.Get("/outstanding", handlers.HandleGetOutstandingBalances(db))
				r.Get("/units/{unitID}/ledger", handlers.HandleGetUnitLedger(db))
				r.Post("/payments", handlers.HandleCreatePayment(db))
				r.Get("/payments", handlers.HandleListPayments(db))
			})

			// Maintenance Module
			r.Route("/maintenance", func(r chi.Router) {
				r.Post("/", handlers.HandleCreateMaintenanceRequest(db))
				r.Get("/", handlers.HandleListMaintenanceRequests(db))
				r.Patch("/{id}/status", handlers.HandleUpdateMaintenanceStatus(db))
				r.With(middleware.RequireRole(db, "ADMIN")).Patch("/{id}/assign", handlers.HandleAssignMaintenance(db))
			})

			// Vendor Module
			r.Route("/vendors", func(r chi.Router) {
				r.Get("/", handlers.HandleListVendors(db))
				r.Post("/", handlers.HandleCreateVendor(db))
				r.Get("/{id}", handlers.HandleGetVendor(db))
				r.Delete("/{id}", handlers.HandleDeleteVendor(db))
			})

			// Announcements Module
			r.Route("/announcements", func(r chi.Router) {
				r.Get("/", annHandler.List)
				r.With(middleware.RequireRole(db, "ADMIN")).Post("/", annHandler.Create)
				r.With(middleware.RequireRole(db, "ADMIN")).Delete("/{id}", annHandler.Delete)
			})

			// Notifications Module
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", handlers.HandleListNotifications(db))
				r.Get("/unread-count", handlers.HandleGetUnreadCount(db))
				r.Put("/mark-all-read", handlers.HandleMarkAllNotificationsRead(db))
				r.Patch("/{id}/read", handlers.HandleMarkNotificationRead(db))
			})
		})
	})

	// Static File Serving (for the Svelte frontend)
	// This assumes your frontend build is in a 'dist' or 'public' folder
	workDir, _ := os.Getwd()

	// Handle static assets (JS, CSS, Images)
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(workDir, "frontend/dist/assets")))))

	// Catch-all: Route everything else to index.html for Svelte SPA routing
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(workDir, "frontend/dist/index.html"))
	})

	return r
}
