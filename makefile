include .env
export

export PROJECT_ROOT=$(shell pwd)

.PHONY: run-all stop-all kitchen-run rest-demo-run env-up env-down env-cleanup migrate-create migrate-up migrate-down migrate-action env-port-forward env-port-close lint

run-all:
	docker compose up -d --build

stop-all:
	docker compose down -v

kitchen-run:
	@export POSTGRES_HOST=localhost && \
	export LOGGER_FOLDER=${PROJECT_ROOT}/out_local/logs && \
	export RESTAURANT_URL=http://localhost:8081 && \
	go run ${PROJECT_ROOT}/cmd/kitchen/main.go

rest-demo-run:
	@go run ${PROJECT_ROOT}/cmd/restaurant-demo/main.go

env-up:
	docker compose up -d kitchen_postgres

env-down:
	docker compose down kitchen_postgres


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq"; \
		exit 1; \
	fi; \
	docker compose run --rm kitchen_postgres_migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@$(MAKE) --no-print-directory migrate-action action=up

migrate-down:
	@$(MAKE) --no-print-directory migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует параметр action"; \
		exit 1; \
	fi; \
	docker compose run --rm kitchen_postgres_migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@kitchen_postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"

env-port-forward:
	@docker compose up -d port_forwarder

env-port-close:
	@docker compose down port_forwarder

lint:
	golangci-lint run