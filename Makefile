include .env
export

export PROJECT_ROOT=$(shell pwd)

# PostgreSQL
env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "WARN: Do you want to delete all volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Volume files cleared"; \
	else \
		echo "Cancelled"; \
	fi

# Socat
env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

# Migrations
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Required parameter seq is missing.\nExample: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	@docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Required parameter 'action' is missing.\nAvailable actions:\nup (count)\ndown (count)\nversion\nforce (count)\nThe 'action' parameter may include 'count'"; \
		exit 1; \
	fi; \

	@docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

#Logs
logs-cleanup:
	@read -p "WARN: Do you want to delete all log files files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Log files cleared"; \
	else \
		echo "Cancelled"; \
	fi

#To Do App
todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

postgres:
	@make env-up && make env-port-forward