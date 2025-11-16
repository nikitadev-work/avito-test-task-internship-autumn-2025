# docker-compose.yml or docker-compose.dev.yml
COMPOSE_FILE = ./ops/docker-compose.dev.yml

.PHONY: dev-up dev-down dev-restart down-volumes \
	dev-logs-pr-manager-service dev-logs-all

dev-up:
	@echo "Starting dev environment..."
	@docker compose -f $(COMPOSE_FILE) up --build

dev-down:
	@echo "Stopping dev environment..."
	docker compose -f $(COMPOSE_FILE) down

dev-restart: 
	@echo "Restarting dev environment..."
	@dev-down dev-up

clear-volumes:
	@echo "Removing docker volumes..."
	@docker compose -f $(COMPOSE_FILE) down -v
	@docker volume prune -f

dev-logs-pr-manager-service:
	@echo "Logs of pr-manager-service:"
	@docker compose -f $(COMPOSE_FILE) logs -f pr-manager-service

dev-logs-all:
	@echo "All logs of pr-manager-service:"
	docker compose -f $(COMPOSE_FILE) logs
