# Resolve air from PATH, falling back to the usual GOPATH install location.
AIR ?= $(shell command -v air 2>/dev/null || echo "$(shell go env GOPATH)/bin/air")
AIR_VERSION ?= v1.52.3
HOST_PORT ?= 8080

.PHONY: help tools dev build run site preview test fmt vet clean \
        docker-build docker-run up up-dev down logs sh

help: ## Show this help
	@grep -hE '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ----------------------------------------------------------------- local dev

tools: ## Install the development tools (air)
	go install github.com/air-verse/air@$(AIR_VERSION)

dev: ## Run the dev server with hot reload
	@test -x "$(AIR)" || { \
		echo "air not found at $(AIR)."; \
		echo "Install it with: make tools"; \
		exit 1; \
	}
	@"$(AIR)" -c .air.toml

build: ## Build the production binary
	go build -o server .

run: build ## Build and run without hot reload
	PORT=$(HOST_PORT) ./server

# ---------------------------------------------------------------------- site

site: ## Build the static site into dist/
	go run ./cmd/build -out dist

preview: site ## Build it and serve dist/ exactly as GitHub Pages will
	@echo "Serving dist/ on http://localhost:$(HOST_PORT)"
	@cd dist && python3 -m http.server $(HOST_PORT)

# --------------------------------------------------------------------- checks

test: ## Run the suite
	go test ./...

fmt: ## Format everything
	gofmt -w .

vet: ## Vet everything
	go vet ./...

# --------------------------------------------------------------------- docker

docker-build: ## Build the production image
	docker build --target runtime -t golang-codes:local .

docker-run: docker-build ## Run the production image directly
	docker run --rm -p $(HOST_PORT):8080 golang-codes:local

up: ## Start the production container via compose
	docker compose up --build -d
	@echo "Serving on http://localhost:$(HOST_PORT)"

up-dev: ## Start the hot-reloading container via compose
	docker compose --profile dev up --build

down: ## Stop and remove the containers
	docker compose --profile dev down

logs: ## Follow the container logs
	docker compose logs -f

sh: ## Shell into the running production container
	docker compose exec web sh

# ---------------------------------------------------------------------- clean

clean: ## Remove build artifacts
	rm -f server
	rm -rf tmp tmp-docker dist
