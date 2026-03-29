.PHONY: up down migrate run test frontend-install frontend

# Start the database container
up:
	docker compose up -d

# Stop the database container
down:
	docker compose down

# Apply the database schema
migrate:
	docker exec -i gated_community_db psql -U user -d gated_db < internal/database/schema.sql

# Run the Go application
run:
	go run cmd/api/main.go

# Run all tests
test:
	go test -v ./...

# Install frontend dependencies
frontend-install:
	cd web && npm install

# Run the SvelteKit development server
frontend:
	cd web && npm run dev

# Full setup: start db, wait a moment, and migrate
setup: up
	@echo "Waiting for database to be ready..."
	@sleep 3
	$(MAKE) migrate