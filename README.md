<div align="center">

# 🚀 gostart

**Enterprise-grade Go project scaffolder**

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/iyashjayesh/gostart/actions/workflows/ci.yml/badge.svg)](https://github.com/iyashjayesh/gostart/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iyashjayesh/gostart)](https://goreportcard.com/report/github.com/iyashjayesh/gostart)

*Think `create-react-app`, but for Go services.*

</div>

---

## What is gostart?

`gostart` generates **opinionated, production-quality Go project boilerplate** via an interactive CLI. Stop copy-pasting skeleton code between projects. Start with:

- ✅ Graceful shutdown with OS signal handling
- ✅ Structured JSON logging (slog, zerolog, or zap)
- ✅ Your choice of HTTP framework (gin, fiber, chi, echo, gorilla/mux)
- ✅ Viper-powered config with `.env` support
- ✅ Optional infra clients: Redis, Kafka, NATS
- ✅ Multi-stage distroless Dockerfile + docker-compose
- ✅ Makefile, GitHub Actions CI, golangci-lint - ready to go on day one

---

## Install

### Go install
```bash
go install github.com/iyashjayesh/gostart@latest
```

### Homebrew (coming soon)
```bash
brew install gostart-dev/tap/gostart
```

### From source
```bash
git clone https://github.com/iyashjayesh/gostart.git
cd gostart
make install
```

---

## Usage

### Interactive mode
```bash
gostart init my-api
```

Sample prompt flow:
```
 ██████╗  ██████╗ ███████╗████████╗ █████╗ ██████╗ ████████╗
...

? Module name: (github.com/your-org/my-api)
? Go version: (1.23)
? HTTP framework: gin
? Structured logger: slog (stdlib)
? Add Viper for config & env management? (Y/n)
? Add Redis client (go-redis)? (y/N)
? Add Kafka client (franz-go)? (y/N)
? Add NATS client? (y/N)
? Add Dockerfile + docker-compose? (Y/n)
? Add Makefile? (Y/n)
? Add GitHub Actions CI? (Y/n)
? Add golangci-lint config? (Y/n)
? Add Swagger/OpenAPI scaffold? (y/N)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Scaffolding project: my-api
  Module:    github.com/your-org/my-api
  Go:        1.23
  Framework: gin
  Logger:    slog
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓ cmd/main.go
  ✓ go.mod
  ✓ .gitignore
  ✓ .env.example
  ✓ internal/config/config.go
  ✓ internal/handler/handler.go
  ✓ internal/server/server.go
  ✓ Dockerfile
  ✓ docker-compose.yml
  ✓ Makefile
  ✓ .github/workflows/ci.yml
  ✓ .golangci.yml
  → Running go mod tidy...
  ✓ go mod tidy

  ✔ Project created successfully!

  Next steps:
    cd my-api
    cp .env.example .env
    make docker-up
    make run
```

### Non-interactive / CI mode
```bash
gostart init my-api --defaults
gostart init my-api --framework fiber --logger zap --redis --kafka --docker
```

---

## Flags

| Flag            | Default   | Description                                       |
|-----------------|-----------|---------------------------------------------------|
| `--module`      | `""`      | Go module path                                    |
| `--go-version`  | `1.23`    | Go version (`1.21`, `1.22`, `1.23`)               |
| `--framework`   | `gin`     | HTTP framework (`gin\|fiber\|chi\|echo\|gorilla\|none`) |
| `--logger`      | `slog`    | Structured logger (`slog\|zerolog\|zap`)           |
| `--viper`       | `true`    | Add Viper for config & env management             |
| `--redis`       | `false`   | Add Redis client (go-redis/v9)                    |
| `--kafka`       | `false`   | Add Kafka client (franz-go)                       |
| `--nats`        | `false`   | Add NATS client                                   |
| `--docker`      | `true`    | Add Dockerfile + docker-compose                   |
| `--makefile`    | `true`    | Add Makefile                                      |
| `--github`      | `true`    | Add GitHub Actions CI                             |
| `--lint`        | `true`    | Add golangci-lint config                          |
| `--swagger`     | `false`   | Add Swagger/OpenAPI scaffold                      |
| `--defaults`    | `false`   | Skip all prompts, use recommended defaults        |
| `--output`      | `.`       | Output directory                                  |

---

## Generated Project Structure

```
my-api/
├── cmd/
│   └── main.go                  # Entrypoint with graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go            # Viper (or stdlib) config loader
│   ├── server/
│   │   └── server.go            # HTTP server for chosen framework
│   └── handler/
│       └── handler.go           # HTTP handlers
├── pkg/
│   ├── redis/redis.go           # (if selected) go-redis wrapper
│   ├── kafka/kafka.go           # (if selected) franz-go producer+consumer
│   └── nats/nats.go             # (if selected) NATS client wrapper
├── .github/
│   └── workflows/ci.yml
├── .env.example
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── Makefile
```

---

## Generated Makefile Targets

| Target             | Description                          |
|--------------------|--------------------------------------|
| `make run`         | `go run ./cmd/main.go`               |
| `make build`       | Build binary to `bin/<name>`         |
| `make test`        | `go test -race -cover ./...`         |
| `make test-coverage` | Generate `coverage.html` report    |
| `make lint`        | `golangci-lint run ./...`            |
| `make fmt`         | `gofmt -s -w . && goimports -w .`    |
| `make tidy`        | `go mod tidy`                        |
| `make docker-up`   | `docker compose up -d`               |
| `make docker-down` | `docker compose down`                |
| `make docker-logs` | `docker compose logs -f app`         |
| `make clean`       | Remove `bin/`, `coverage.*`          |
| `make install-tools` | Install golangci-lint + goimports  |

---

## Supported Frameworks

| Flag value   | Import path                        |
|--------------|------------------------------------|
| `gin`        | `github.com/gin-gonic/gin`         |
| `fiber`      | `github.com/gofiber/fiber/v2`      |
| `chi`        | `github.com/go-chi/chi/v5`         |
| `echo`       | `github.com/labstack/echo/v4`      |
| `gorilla`    | `github.com/gorilla/mux`           |
| `none`       | stdlib `net/http`                  |

---

## Philosophy

**Idiomatic Go.** Generated code uses standard patterns: `context` propagation, structured logging through `slog` as default, `os.Exit(1)` on unrecoverable errors, `signal.NotifyContext` for graceful shutdown.

**Minimal opinions.** gostart picks sensible defaults but lets you override everything. Choose your framework, choose your logger, opt-in or out of every infrastructure component.

**Production defaults from day one.** Every generated project includes: server timeouts, race-detected tests, a multi-stage distroless Docker build, and a CI pipeline that runs lint before merge.

---

## Contributing

```bash
# Clone and build
git clone https://github.com/iyashjayesh/gostart.git
cd gostart
make build

# Run the smoke test
make smoke-test
```

PRs welcome! Please open an issue first for major changes.

---

## License

[MIT](LICENSE) © gostart contributors
