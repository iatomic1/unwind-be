GO = go
DB_CONTAINER = unwind-be-pg
DB_PASSWORD = password
DB_PORT = 5431
DB_URL = postgres://postgres:$(DB_PASSWORD)@localhost:$(DB_PORT)/postgres?sslmode=disable
MIGRATION_DIR = internal/db/migrations
DOCS_DIR = internal/docs
SWAGGER_ENTRY = cmd/api/main.go
SWAGGER_FLAGS = --parseDepth 1 --parseDependency --parseInternal --o $(DOCS_DIR)

.PHONY: restart clean migrate docs dev

clean:
	rm -rf ./tmp

migrate:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

convert-docs:
	swagger2openapi $(DOCS_DIR)/swagger.json > $(DOCS_DIR)/openapi.json

add-prefix:
	python3 ./scripts/add_prefix.py

docs:
	@rm -rf ./internal/docs/
	swag init -g $(SWAGGER_ENTRY) $(SWAGGER_FLAGS)
	@sleep 1
	$(MAKE) convert-docs
	# $(MAKE) add-prefix

create:
	@goose create $(n) sql -dir $(MIGRATION_DIR)

stop-container:
	@echo "Stopping existing container..."
	docker stop $(DB_CONTAINER) || true


remove-container:
	@echo "Removing container..."
	docker rm $(DB_CONTAINER) || true

start-container:
	@echo "Starting new container..."
	docker run --name $(DB_CONTAINER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-p $(DB_PORT):5432 \
		-d postgres

restart-db:
	$(MAKE) stop-container
	$(MAKE) remove-container
	$(MAKE) start-container
	@echo "Waiting for PostgreSQL to start..."
	@sleep 3
	@echo "Enabling UUID extension..."
	@docker exec -i $(DB_CONTAINER) psql -U postgres -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'

psql:
	@docker exec -it $(DB_CONTAINER) psql -U postgres

gen:
	@python3 ./scripts/add_tags.py
	@sqlc generate

dev:
	$(MAKE) clean
	docker start $(DB_CONTAINER)
	$(MAKE) migrate
	$(MAKE) gen
	$(MAKE) docs
	air
