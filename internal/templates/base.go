package templates

import "github.com/iyashjayesh/goscaf/internal/config"

// MainGo is the template for cmd/main.go
const MainGo = `{{if .Swagger}}// @title           {{.ProjectName}} API
// @version         1.0.0
// @description     API documentation for {{.ProjectName}}
// @host            localhost:8080
// @BasePath        /api/v1
{{end}}package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/server"
)

func main() {
	// Set up structured JSON logger as default
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create and start server
	srv := server.New(cfg, logger)
	if err := srv.Start(ctx); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
`

// GoMod is the template for go.mod
const GoMod = `module {{.ModuleName}}

go {{.GoVersion}}
{{if ne .Framework "none"}}
require (
	{{.FrameworkImport}} latest
{{- if eq .Logger "zerolog"}}
	github.com/rs/zerolog v1.32.0
{{- end}}
{{- if eq .Logger "zap"}}
	go.uber.org/zap v1.27.0
{{- end}}
{{- if .Viper}}
	github.com/spf13/viper v1.18.2
{{- end}}
{{- if .Redis}}
	github.com/redis/go-redis/v9 v9.5.1
{{- end}}
{{- if .Kafka}}
	github.com/twmb/franz-go v1.16.1
{{- end}}
{{- if .NATS}}
	github.com/nats-io/nats.go v1.33.1
{{- end}}
)
{{else}}
require (
{{- if eq .Logger "zerolog"}}
	github.com/rs/zerolog v1.32.0
{{- end}}
{{- if eq .Logger "zap"}}
	go.uber.org/zap v1.27.0
{{- end}}
{{- if .Viper}}
	github.com/spf13/viper v1.18.2
{{- end}}
{{- if .Redis}}
	github.com/redis/go-redis/v9 v9.5.1
{{- end}}
{{- if .Kafka}}
	github.com/twmb/franz-go v1.16.1
{{- end}}
{{- if .NATS}}
	github.com/nats-io/nats.go v1.33.1
{{- end}}
)
{{end}}
`

// GitIgnore is the template for .gitignore
const GitIgnore = `# Binaries
/bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out
coverage.out
coverage.html

# Env files
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Vendor
vendor/
`

// EnvExample is the template for .env.example
const EnvExample = `# Application
APP_NAME={{.ProjectName}}
APP_ENV=development
APP_PORT=8080
{{if .Redis}}
# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
{{end}}{{if .Kafka}}
# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID={{.ProjectName}}-consumer
{{end}}{{if .NATS}}
# NATS
NATS_URL=nats://localhost:4222
NATS_NAME={{.ProjectName}}
{{end}}
`

// ConfigGo returns the internal/config/config.go template (viper or stdlib).
func ConfigGo(useViper bool) string {
	if useViper {
		return configGoViper
	}
	return configGoStdlib
}

const configGoViper = `package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App   AppConfig
{{- if .Redis}}
	Redis RedisConfig
{{- end}}
{{- if .Kafka}}
	Kafka KafkaConfig
{{- end}}
{{- if .NATS}}
	NATS  NATSConfig
{{- end}}
}

// AppConfig holds app-level configuration.
type AppConfig struct {
	Name string
	Env  string
	Port int
}
{{if .Redis}}
// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
{{end}}{{if .Kafka}}
// KafkaConfig holds Kafka connection configuration.
type KafkaConfig struct {
	Brokers []string
	GroupID string
}
{{end}}{{if .NATS}}
// NATSConfig holds NATS connection configuration.
type NATSConfig struct {
	URL  string
	Name string
}
{{end}}
// Load reads configuration from environment variables and optional .env file.
func Load() (*Config, error) {
	v := viper.New()

	// Read .env file (not fatal if missing)
	v.SetConfigFile(".env")
	v.SetConfigType("dotenv")
	_ = v.ReadInConfig()

	// Environment variable settings
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("app.name", "{{.ProjectName}}")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
{{if .Redis}}
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
{{end}}{{if .Kafka}}
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.group_id", "{{.ProjectName}}-consumer")
{{end}}{{if .NATS}}
	v.SetDefault("nats.url", "nats://localhost:4222")
	v.SetDefault("nats.name", "{{.ProjectName}}")
{{end}}

	cfg := &Config{
		App: AppConfig{
			Name: v.GetString("APP_NAME"),
			Env:  v.GetString("APP_ENV"),
			Port: v.GetInt("APP_PORT"),
		},
{{if .Redis}}
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
{{end}}{{if .Kafka}}
		Kafka: KafkaConfig{
			Brokers: v.GetStringSlice("KAFKA_BROKERS"),
			GroupID: v.GetString("KAFKA_GROUP_ID"),
		},
{{end}}{{if .NATS}}
		NATS: NATSConfig{
			URL:  v.GetString("NATS_URL"),
			Name: v.GetString("NATS_NAME"),
		},
{{end}}
	}

	// Apply viper defaults if env vars are empty
	if cfg.App.Name == "" {
		cfg.App.Name = v.GetString("app.name")
	}
	if cfg.App.Env == "" {
		cfg.App.Env = v.GetString("app.env")
	}
	if cfg.App.Port == 0 {
		cfg.App.Port = v.GetInt("app.port")
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.App.Port < 1 || c.App.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", c.App.Port)
	}
	return nil
}
`

const configGoStdlib = `package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	App   AppConfig
{{- if .Redis}}
	Redis RedisConfig
{{- end}}
{{- if .Kafka}}
	Kafka KafkaConfig
{{- end}}
{{- if .NATS}}
	NATS  NATSConfig
{{- end}}
}

// AppConfig holds app-level configuration.
type AppConfig struct {
	Name string
	Env  string
	Port int
}
{{if .Redis}}
// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
{{end}}{{if .Kafka}}
// KafkaConfig holds Kafka connection configuration.
type KafkaConfig struct {
	Brokers []string
	GroupID string
}
{{end}}{{if .NATS}}
// NATSConfig holds NATS connection configuration.
type NATSConfig struct {
	URL  string
	Name string
}
{{end}}
// Load reads configuration from environment variables.
func Load() (*Config, error) {
	port := 8080
	if p := os.Getenv("APP_PORT"); p != "" {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid APP_PORT: %w", err)
		}
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "{{.ProjectName}}"),
			Env:  getEnv("APP_ENV", "development"),
			Port: port,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func (c *Config) validate() error {
	if c.App.Port < 1 || c.App.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", c.App.Port)
	}
	return nil
}
`

// HandlerGo is the template for internal/handler/handler.go
const HandlerGo = `package handler

import (
{{- if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
{{- else if eq .Framework "fiber"}}
	"github.com/gofiber/fiber/v2"
{{- else if eq .Framework "chi"}}
	"encoding/json"
	"net/http"
{{- else if eq .Framework "echo"}}
	"github.com/labstack/echo/v4"
	"net/http"
{{- else if eq .Framework "gorilla"}}
	"encoding/json"
	"net/http"
{{- else}}
	"encoding/json"
	"net/http"
{{- end}}

	"{{.ModuleName}}/internal/config"
)

// Handler holds application dependencies.
type Handler struct {
	cfg *config.Config
}

// New creates a new Handler.
func New(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// healthResponse is the health check response body.
type healthResponse struct {
	Status string ` + "`" + `json:"status"` + "`" + `
	App    string ` + "`" + `json:"app"` + "`" + `
	Env    string ` + "`" + `json:"env"` + "`" + `
}
{{if eq .Framework "gin"}}
{{- if .Swagger}}
// Health godoc
// @Summary     Health check
// @Description Returns service health status
// @Tags        health
// @Produce     json
// @Success     200  {object}  healthResponse
// @Router      /health [get]
{{- end}}
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, healthResponse{
		Status: "ok",
		App:    h.cfg.App.Name,
		Env:    h.cfg.App.Env,
	})
}
{{else if eq .Framework "fiber"}}
{{- if .Swagger}}
// Health godoc
// @Summary     Health check
// @Description Returns service health status
// @Tags        health
// @Produce     json
// @Success     200  {object}  healthResponse
// @Router      /health [get]
{{- end}}
func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(healthResponse{
		Status: "ok",
		App:    h.cfg.App.Name,
		Env:    h.cfg.App.Env,
	})
}
{{else if eq .Framework "echo"}}
{{- if .Swagger}}
// Health godoc
// @Summary     Health check
// @Description Returns service health status
// @Tags        health
// @Produce     json
// @Success     200  {object}  healthResponse
// @Router      /health [get]
{{- end}}
func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, healthResponse{
		Status: "ok",
		App:    h.cfg.App.Name,
		Env:    h.cfg.App.Env,
	})
}
{{else}}
{{- if .Swagger}}
// Health godoc
// @Summary     Health check
// @Description Returns service health status
// @Tags        health
// @Produce     json
// @Success     200  {object}  healthResponse
// @Router      /health [get]
{{- end}}
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{
		Status: "ok",
		App:    h.cfg.App.Name,
		Env:    h.cfg.App.Env,
	})
}
{{end}}
`

// Ensure config package is imported in handler template
var _ = config.ProjectConfig{}
