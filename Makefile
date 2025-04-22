# Makefile for microservices project

# Variables
DOCKER_COMPOSE = docker-compose
GO = go
GOFLAGS = -v
GOCMD = $(GO) $(GOFLAGS)
PRODUCT_SERVICE_DIR = services/product-service
POSTGRES_CONTAINER = postgres
DB_NAME = products
DB_USER = postgres

# Docker commands
.PHONY: docker-up
docker-up:
    $(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
    $(DOCKER_COMPOSE) down

.PHONY: docker-logs
docker-logs:
    $(DOCKER_COMPOSE) logs -f

# Build commands
.PHONY: build-all
build-all: build-product-service

.PHONY: build-product-service
build-product-service:
    $(GOCMD) build -o bin/product-service $(PRODUCT_SERVICE_DIR)/cmd/api/main.go

# Run commands
.PHONY: run-product-service
run-product-service:
    $(GOCMD) run $(PRODUCT_SERVICE_DIR)/cmd/api/main.go

# Database commands
.PHONY: db-seed-categories
db-seed-categories:
    docker exec -i $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < scripts/db/seed/categories.sql

.PHONY: db-seed-products
db-seed-products:
    docker exec -i $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < scripts/db/seed/products.sql

.PHONY: db-seed
db-seed: db-seed-categories db-seed-products

.PHONY: db-reset
db-reset:
    docker exec -i $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "TRUNCATE products, categories CASCADE;"

.PHONY: db-query-products
db-query-products:
    docker exec -i $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT id, name, price, sku, category_id FROM products LIMIT 10;"

.PHONY: db-query-categories
db-query-categories:
    docker exec -i $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT id, name FROM categories;"

# Test commands
.PHONY: test-all
test-all:
    $(GOCMD) test ./...

.PHONY: test-product-service
test-product-service:
    $(GOCMD) test ./$(PRODUCT_SERVICE_DIR)/...

# Utility commands
.PHONY: clean
clean:
    rm -rf bin/
    rm -rf tmp/

.PHONY: help
help:
    @echo "Available commands:"
    @echo "  docker-up               - Start all containers"
    @echo "  docker-down             - Stop all containers"
    @echo "  docker-logs             - Show container logs"
    @echo "  build-all               - Build all services"
    @echo "  build-product-service   - Build product service"
    @echo "  run-product-service     - Run product service"
    @echo "  db-seed-categories      - Seed categories data"
    @echo "  db-seed-products        - Seed products data"
    @echo "  db-seed                 - Seed all data"
    @echo "  db-reset                - Reset database (delete all data)"
    @echo "  db-query-products       - Query products table"
    @echo "  db-query-categories     - Query categories table"
    @echo "  test-all                - Run all tests"
    @echo "  test-product-service    - Run product service tests"
    @echo "  clean                   - Clean build artifacts"