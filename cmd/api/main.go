package main

// a comment
import (
	"fmt"
	"log"
	"net/http"

	"gated-community-api/config"
	"gated-community-api/internal/database"
	"gated-community-api/internal/routes"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Database
	dbPool := database.NewPostgresPool(cfg.DatabaseURL)
	defer dbPool.Close()

	// 3. Setup Router
	router := routes.ConfigureRoutes(dbPool)

	// 4. Start Server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s in %s mode", serverAddr, cfg.Env)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
