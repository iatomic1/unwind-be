GO = go
DB_CONTAINER = unwind-be-pg
DB_PASSWORD = password
DB_PORT = 5431
DB_URL = postgres://postgres:$(DB_PASSWORD)@localhost:$(DB_PORT)/postgres?sslmode=disable
MIGRATION_DIR = internal/db/migrations
SWAGGER_ENTRY = cmd/api/main.go
SWAGGER_FLAGS = --parseDepth 1 --parseDependency --parseInternal

.PHONY: restart clean migrate docs dev

clean:
	rm -rf ./tmp

migrate:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

docs:
	rm -rf ./docs/
	swag init -g $(SWAGGER_ENTRY) $(SWAGGER_FLAGS)

restart-db:
	@echo "Stopping existing container..."
	docker stop $(DB_CONTAINER) || true
	@echo "Removing container..."
	docker rm $(DB_CONTAINER) || true
	@echo "Starting new container..."
	docker run --name $(DB_CONTAINER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-p $(DB_PORT):5432 \
		-d postgres
	@echo "Waiting for PostgreSQL to start..."
	@sleep 3
	@echo "Enabling UUID extension..."
	@docker exec -i $(DB_CONTAINER) psql -U postgres -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'

dev:
	docker start $(DB_CONTAINER)
	rm -rf ./tmp ./docs/
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up
	sleep 1
	python3 new.py
	sqlc generate
	sleep 1
	swag init -g $(SWAGGER_ENTRY) $(SWAGGER_FLAGS)
	air
