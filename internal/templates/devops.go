package templates

// DockerfileTemplate is the template for the Dockerfile.
const DockerfileTemplate = `# ── Build Stage ──────────────────────────────────────────────
FROM golang:{{.GoVersion}}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/{{.ProjectName}} \
    ./cmd/main.go

# ── Runtime Stage ─────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/bin/{{.ProjectName}} /app/{{.ProjectName}}

EXPOSE 8080

ENTRYPOINT ["/app/{{.ProjectName}}"]
`

// DockerCompose is the template for docker-compose.yml.
const DockerCompose = `version: "3.9"

services:
  app:
    build: .
    env_file: .env
    ports:
      - "8080:8080"
    restart: unless-stopped
{{- if or .Redis .Kafka}}
    depends_on:
{{- if .Redis}}
      redis:
        condition: service_healthy
{{- end}}
{{- if .Kafka}}
      kafka:
        condition: service_healthy
{{- end}}
{{- end}}
{{if .Redis}}
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
{{end}}{{if .Kafka}}
  zookeeper:
    image: confluentinc/cp-zookeeper:7.6.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    restart: unless-stopped

  kafka:
    image: confluentinc/cp-kafka:7.6.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "kafka-topics --bootstrap-server localhost:9092 --list"]
      interval: 30s
      timeout: 10s
      retries: 5
{{end}}{{if .NATS}}
  nats:
    image: nats:2-alpine
    command: ["--http_port", "8222"]
    ports:
      - "4222:4222"
      - "8222:8222"
    restart: unless-stopped
{{end}}
`

// MakefileTemplate is the template for the Makefile.
const MakefileTemplate = `APP_NAME := {{.ProjectName}}
BUILD_DIR := bin
BINARY := $(BUILD_DIR)/$(APP_NAME)

.PHONY: run build test test-coverage lint fmt tidy \
        docker-up docker-down docker-logs clean install-tools

run: ## Run the application
	go run ./cmd/main.go

build: ## Build binary to bin/$(APP_NAME)
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BINARY) ./cmd/main.go

test: ## Run tests with race detection and coverage
	go test -race -cover ./...

test-coverage: ## Generate HTML coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format code
	gofmt -s -w .
	goimports -w .

tidy: ## Run go mod tidy
	go mod tidy

docker-up: ## Start all services
	docker compose up -d

docker-down: ## Stop all services
	docker compose down

docker-logs: ## Follow app logs
	docker compose logs -f app

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)/ coverage.out coverage.html

install-tools: ## Install dev tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
`

// GitHubCI is the template for .github/workflows/ci.yml.
const GitHubCI = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "{{.GoVersion}}"
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run go vet
        run: go vet ./...

      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "{{.GoVersion}}"
          cache: true

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "{{.GoVersion}}"
          cache: true

      - name: Build
        run: CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/{{.ProjectName}} ./cmd/main.go
`

// GolangCI is the template for .golangci.yml.
const GolangCI = `version: "2"
linters:
  enable:
    - bodyclose
    - contextcheck
    - gocritic
    - misspell
    - noctx
    - revive
    - wrapcheck
  settings:
    revive:
      rules:
        - name: exported
        - name: unused-parameter
  exclusions:
    generated: lax
    presets:
      - comments
      - common-false-positives
      - legacy
      - std-error-handling
    rules:
      - linters:
          - errcheck
          - wrapcheck
        path: _test\.go
    paths:
      - third_party$
      - builtin$
      - examples$
formatters:
  enable:
    - gofmt
    - goimports
  settings:
    gofmt:
      simplify: true
    goimports:
      local-prefixes:
        - {{.ModuleName}}
  exclusions:
    generated: lax
    paths:
      - third_party$
      - builtin$
      - examples$
`
