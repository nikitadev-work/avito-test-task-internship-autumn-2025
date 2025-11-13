dev-up:

dev-down:

run-swagger:
	docker run --rm -p 8082:8080 \
	-e SWAGGER_JSON=/api/openapi.yml \
	-v ./docs/contracts/pr-manager-service-openapi.yml:/api/openapi.yml \
	swaggerapi/swagger-ui