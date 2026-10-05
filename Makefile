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

tgbot-set-webhook:
	@curl -X POST "https://api.telegram.org/bot$(TGBOT_TOKEN)/setWebhook" \
		-H "Content-Type: application/json" \
		-d '{"url": "$(BASE_URL)/api/telegram/webhook", "secret_token": "$(TGBOT_SERVER_WEBHOOK_SECRET)"}'

tgbot-webhook-info:
	@curl -X GET "https://api.telegram.org/bot$(TGBOT_TOKEN)/getWebhookInfo"

tgbot-delete-webhook:
	@curl -X POST "https://api.telegram.org/bot$(TGBOT_TOKEN)/deleteWebhook" 

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
	@rm -f $(PROJECT_ROOT)/backend/tgbot/.env
	@add() { [ -n "$$2" ] && echo "$$1=$$2" >> $(PROJECT_ROOT)/backend/tgbot/.env || true; }; \
	add TGBOT_CLIENT_BOT_TOKEN "$(TGBOT_TOKEN)"; \
	add TGBOT_CLIENT_API_URL "$(TGBOT_API_URL)"; \
	add TGBOT_CLIENT_TIMEOUT "$(TGBOT_CLIENT_TIMEOUT)"; \
	add TGBOT_SERVER_WEBHOOK_SECRET "$(TGBOT_SERVER_WEBHOOK_SECRET)"; \
	add TGBOT_SERVER_WORKERS "$(TGBOT_SERVER_WORKERS)"; \
	add TGBOT_SERVER_QUEUE_SIZE "$(TGBOT_SERVER_QUEUE_SIZE)"; \
	add LOG_LEVEL "$(TGBOT_LOG_LEVEL)"; \
	add LOG_FOLDER "$(TGBOT_LOG_FOLDER)"; \
	add HTTP_SERVER_ADDR "$(TGBOT_SERVER_HTTP_ADDR)"; \
	add HTTP_SERVER_PORT "$(TGBOT_SERVER_HTTP_PORT)"; \
	add HTTP_SERVER_SHUTDOWN_TIMEOUT "$(TGBOT_SERVER_HTTP_SHUTDOWN_TIMEOUT)"

run-main: generate-envs
	@cd $(PROJECT_ROOT)/backend/main && \
	set -a; . ./.env; set +a; \
	go run ./cmd/main/main.go

run-tgbot-server: generate-envs
	@cd $(PROJECT_ROOT)/backend/tgbot && \
	set -a; . ./.env; set +a; \
	go run ./cmd/tgbot/main.go
