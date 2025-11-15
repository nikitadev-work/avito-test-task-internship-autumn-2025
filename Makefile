# docker-compose.yml or docker-compose.dev.yml
COMPOSE_FILE = ./ops/docker-compose.dev.yml

.PHONY: dev-up dev-down dev-restart dev-logs run-swagger

dev-up:
	docker compose -f $(COMPOSE_FILE) up --build

dev-down:
	docker compose -f $(COMPOSE_FILE) down

dev-restart: dev-down dev-up

dev-logs-pr-manager-service:
	docker compose -f $(COMPOSE_FILE) logs -f pr-manager-service

dev-logs-all:
	docker compose -f $(COMPOSE_FILE) logs

run-swagger:
	docker run --rm -p 8082:8080 \
		-e SWAGGER_JSON=/api/openapi.yml \
		-v ./docs/contracts/pr-manager-service-openapi.yml:/api/openapi.yml \
		swaggerapi/swagger-ui
