.PHONY: postgres-pull postgres-start postgres-stop postgres-restart db-setup db-reset run test clean

POSTGRES_CONTAINER := events-postgres
POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
POSTGRES_DB := events_db
POSTGRES_PORT := 5432

postgres-pull:
	@echo "Pulling Postgres image..."
	docker pull postgres:15

postgres-start:
	@echo "Starting Postgres container..."
	@docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true
	@docker rm $(POSTGRES_CONTAINER) 2>/dev/null || true
	docker run -d \
		--name $(POSTGRES_CONTAINER) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-p $(POSTGRES_PORT):5432 \
		postgres:15
	@echo "Postgres is running on port $(POSTGRES_PORT)"

postgres-stop:
	@echo "Stopping Postgres container..."
	docker stop $(POSTGRES_CONTAINER)

postgres-restart: postgres-stop postgres-start

db-setup:
	@echo "Creating database..."
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DB);" || true
	@echo "Running migrations..."
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < internal/migrations/001_create_events_table.sql
	@echo "Database setup complete!"

db-reset:
	@echo "Resetting database..."
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -c "DROP DATABASE IF EXISTS $(POSTGRES_DB);"
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DB);"
	@echo "Running migrations..."
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < internal/migrations/001_create_events_table.sql
	@echo "Database reset complete!"

run:
	@echo "Starting application..."
	go run ./cmd/api

test:
	@echo "Running tests..."
	go test ./... -v

clean:
	@echo "Cleaning up..."
	docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true
	docker rm $(POSTGRES_CONTAINER) 2>/dev/null || true
	@echo "Cleanup complete!"
