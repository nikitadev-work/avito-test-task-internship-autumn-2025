# docker-compose.yml or docker-compose.dev.yml
COMPOSE_FILE = ./ops/docker-compose.dev.yml

.PHONY: dev-up dev-down dev-restart clear-volumes \
	dev-logs-pr-manager-service dev-logs-all \
	lint-pr-manager-service lint-common lint test-integration \
	load-create-pr load-reassign load-get-reviews


# Docker compose

dev-up:
	@echo "Starting dev environment..."
	@docker compose -f $(COMPOSE_FILE) up --build

dev-down:
	@echo "Stopping dev environment..."
	@docker compose -f $(COMPOSE_FILE) down

dev-restart: dev-down dev-up
	@echo "Restarting dev environment..."

clear-volumes:
	@echo "Removing docker volumes..."
	@docker compose -f $(COMPOSE_FILE) down -v
	@docker volume prune -f


# Logs

dev-logs-pr-manager-service:
	@echo "Logs of pr-manager-service:"
	@docker compose -f $(COMPOSE_FILE) logs -f pr-manager-service

dev-logs-all:
	@echo "All logs of all services:"
	@docker compose -f $(COMPOSE_FILE) logs


# Linting

lint-pr-manager-service:
	@echo "Linting pr-manager-service package..."
	@cd pr-manager-service && golangci-lint run ./...

lint-common:
	@echo "Linting common package..."
	@cd common/kit && golangci-lint run ./...

lint: lint-pr-manager-service lint-common
	@echo "All lint checks passed."


# Load testing

load-create-pr:
	k6 run ops/load-testing/k6_create_pr.js

load-reassign:
	k6 run ops/load-testing/k6_reassign_pr.js

load-get-reviews:
	k6 run ops/load-testing/k6_get_reviews.js


# Integration tests

test-integration:
	@echo "Running integration tests (service must be running)..."
	@cd pr-manager-service && go test ./internal/integration-tests -count=1
