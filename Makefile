APP_NAME := url-shortener
BIN_DIR ?= bin
BIN ?= $(BIN_DIR)/api
GO ?= go
DOCKER_COMPOSE ?= docker compose
AIR ?= $(GO) run github.com/air-verse/air@latest

GOCACHE ?= /tmp/go-build
HTTP_ADDR ?= :8080
DATABASE_URL ?= postgres://url_shortener:url_shortener@localhost:5432/url_shortener?sslmode=disable
REDIS_ADDR ?= localhost:6379
K6_RATE ?= 2
K6_DURATION ?= 1m
K6_PRE_ALLOCATED_VUS ?= 4
K6_MAX_VUS ?= 8

.PHONY: help fmt tidy test build run dev dev-local up down load-test install-air clean infra-up infra-down api-up api-down logs ps compose-config db-shell redis-cli

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

fmt: ## Format Go files
	$(GO) fmt ./...

tidy: ## Tidy Go modules
	GOCACHE=$(GOCACHE) $(GO) mod tidy

test: ## Run Go tests
	GOCACHE=$(GOCACHE) $(GO) test ./...

build: ## Build the API binary
	mkdir -p $(BIN_DIR)
	GOCACHE=$(GOCACHE) CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/api

run: ## Run the API locally
	GOCACHE=$(GOCACHE) HTTP_ADDR="$(HTTP_ADDR)" DATABASE_URL="$(DATABASE_URL)" REDIS_ADDR="$(REDIS_ADDR)" $(GO) run ./cmd/api

dev: ## Start the full dev stack with Air, observability, Postgres, and Redis
	$(DOCKER_COMPOSE) --profile dev up --build

up: dev ## Alias for dev

load-test: ## Run a safe k6 load test against the dev stack
	K6_RATE="$(K6_RATE)" K6_DURATION="$(K6_DURATION)" K6_PRE_ALLOCATED_VUS="$(K6_PRE_ALLOCATED_VUS)" K6_MAX_VUS="$(K6_MAX_VUS)" $(DOCKER_COMPOSE) --profile dev up -d
	K6_RATE="$(K6_RATE)" K6_DURATION="$(K6_DURATION)" K6_PRE_ALLOCATED_VUS="$(K6_PRE_ALLOCATED_VUS)" K6_MAX_VUS="$(K6_MAX_VUS)" $(DOCKER_COMPOSE) --profile dev --profile loadtest run --rm k6

dev-local: infra-up ## Run the API locally with Air live reload
	GOCACHE=$(GOCACHE) HTTP_ADDR="$(HTTP_ADDR)" DATABASE_URL="$(DATABASE_URL)" REDIS_ADDR="$(REDIS_ADDR)" $(AIR) -c .air.toml

install-air: ## Install Air globally into GOPATH/bin
	$(GO) install github.com/air-verse/air@latest

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

infra-up: ## Start Postgres and Redis
	$(DOCKER_COMPOSE) up -d postgres redis

infra-down: ## Stop and remove Compose services
	$(DOCKER_COMPOSE) down

down: ## Stop and remove Compose services
	$(DOCKER_COMPOSE) down --remove-orphans

api-up: ## Start Postgres, Redis, observability, and the API container
	$(DOCKER_COMPOSE) --profile api up --build

api-down: ## Stop and remove Compose services including API
	$(DOCKER_COMPOSE) --profile api down

logs: ## Follow Compose logs
	$(DOCKER_COMPOSE) logs -f

ps: ## Show Compose service status
	$(DOCKER_COMPOSE) ps

compose-config: ## Validate the dev and load test Compose configuration
	$(DOCKER_COMPOSE) --profile dev --profile loadtest config

db-shell: ## Open a psql shell in Postgres
	$(DOCKER_COMPOSE) exec postgres psql -U url_shortener -d url_shortener

redis-cli: ## Open a Redis CLI shell
	$(DOCKER_COMPOSE) exec redis redis-cli
