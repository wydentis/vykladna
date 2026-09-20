include .env
export

export PROJECT_ROOT=${shell pwd}

export UID := $(shell id -u)
export GID := $(shell id -g)

env-up:
	@docker compose up -d main-postgres

env-down:
	@docker compose down main-postgres

env-reset:
	@docker compose down -v main-postgres

env-restart:
	@docker compose down main-postgres && \
	docker compose up -d main-postgres

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq (aka name) parameter is required"; \
		exit 1; \
	fi; \
	if [ -z "$(db)" ]; then \
		echo "db parameter is required"; \
		exit 1; \
	elif [ "$(db)" = "main" ]; then \
		SERVICE="main-postgres-migrate"; \
	else \
		echo "unknown db: $(db)"; \
		exit 1; \
	fi; \
	mkdir -p $(PROJECT_ROOT)/migrations/$(db); \
	docker compose run --rm $$SERVICE \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action parameter is required"; \
		exit 1; \
	fi; \
	if [ -z "$(db)" ]; then \
		echo "db parameter is required"; \
		exit 1; \
	fi; \
	if [ "$(db)" = "main" ]; then \
		USER="$(MAIN_POSTGRES_USER)"; PASS="$(MAIN_POSTGRES_PASSWORD)"; DBNAME="$(MAIN_POSTGRES_DBNAME)" \
		HOST="main-postgres"; SERVICE="main-postgres-migrate"; \
	else \
		echo "unknown db: $(db)"; \
		exit 1; \
	fi; \
	docker compose run --rm $$SERVICE \
		-path /migrations \
		-database postgres://$$USER:$$PASS@$$HOST:5432/$$DBNAME?sslmode=disable \
		$(action)

generate-envs:
	@rm -f $(PROJECT_ROOT)/backend/main/.env
	@add() { [ -n "$$2" ] && echo "$$1=$$2" >> $(PROJECT_ROOT)/backend/main/.env || true; }; \
	add LOG_LEVEL "$(MAIN_APP_LOG_LEVEL)"; \
	add LOG_FOLDER "$(MAIN_APP_LOG_FOLDER)"; \
	add POSTGRES_HOST "$(MAIN_POSTGRES_HOST)"; \
	add POSTGRES_PORT "$(MAIN_POSTGRES_PORT)"; \
	add POSTGRES_USER "$(MAIN_POSTGRES_USER)"; \
	add POSTGRES_PASSWORD "$(MAIN_POSTGRES_PASSWORD)"; \
	add POSTGRES_DBNAME "$(MAIN_POSTGRES_DBNAME)"; \
	add POSTGRES_TIMEOUT "$(MAIN_POSTGRES_TIMEOUT)"; \
	add HTTP_SERVER_ADDR "$(MAIN_APP_HTTP_ADDR)"; \
	add HTTP_SERVER_PORT "$(MAIN_APP_HTTP_PORT)"; \
	add HTTP_SERVER_SHUTDOWN_TIMEOUT "$(MAIN_APP_HTTP_SHUTDOWN_TIMEOUT)"

run-main: generate-envs
	@cd $(PROJECT_ROOT)/backend/main && \
	set -a; . ./.env; set +a; \
	go run ./cmd/main/main.go