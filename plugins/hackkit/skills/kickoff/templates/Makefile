# AI-buddy — common commands. `make help` lists them.
SHELL := bash
GO_DIR ?= backend
FRONT_DIR ?= frontend
COVERAGE_MIN ?= 85
export COVERAGE_MIN
# First golangci-lint v2 found: on PATH, else in $(go env GOPATH)/bin (an old v1 on PATH would reject the v2 config).
GOLANGCI ?= $(shell for b in golangci-lint "$$(go env GOPATH)/bin/golangci-lint"; do "$$b" version 2>/dev/null | grep -q 'version 2\.' && { echo "$$b"; break; }; done)

.PHONY: help up down logs test test-go test-fe lint cover-check e2e check

help: ## Show targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*## "}{printf "  %-12s %s\n",$$1,$$2}'

up: ## Build and start everything (docker compose, no .env needed), then print URLs
	docker compose up -d --build --wait
	@echo "Frontend: http://localhost:5173   Backend: http://localhost:8080/healthz"

down: ## Stop everything
	docker compose down

logs: ## Follow logs
	docker compose logs -f --tail=100

test-go: ## Go tests with race detector
	cd $(GO_DIR) && go test ./... -race -count=1

test-fe: ## Frontend unit/component tests
	cd $(FRONT_DIR) && npx vitest run

test: test-go test-fe ## All unit tests

lint: ## golangci-lint v2 (go vet if v2 is not installed) + tsc + eslint (if configured)
	@cd $(GO_DIR) && if [ -n "$(GOLANGCI)" ]; then $(GOLANGCI) run ./...; \
	else echo "golangci-lint v2 not found (go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest), running go vet"; go vet ./...; fi
	cd $(FRONT_DIR) && npx tsc --noEmit -p .
	@cd $(FRONT_DIR) && if npm run | grep -q ' lint'; then npm run lint; fi

cover-check: ## Coverage gate (> 80% required; COVERAGE_MIN, default 85), writes coverage/REPORT.md
	bash scripts/coverage-check.sh

e2e: ## Playwright e2e against the running stack (make up first)
	cd $(FRONT_DIR) && npx playwright test

check: lint test cover-check ## Everything that must be green before "done"
