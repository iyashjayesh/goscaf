package config

import "fmt"

// Framework represents the HTTP framework choice.
type Framework string

const (
	FrameworkGin     Framework = "gin"
	FrameworkFiber   Framework = "fiber"
	FrameworkChi     Framework = "chi"
	FrameworkEcho    Framework = "echo"
	FrameworkGorilla Framework = "gorilla"
	FrameworkNone    Framework = "none"
)

// Logger represents the structured logger choice.
type Logger string

const (
	LoggerSlog    Logger = "slog"
	LoggerZerolog Logger = "zerolog"
	LoggerZap     Logger = "zap"
)

// ProjectConfig holds all configuration for a new Go project.
type ProjectConfig struct {
	// Core
	ProjectName string
	ModuleName  string
	GoVersion   string
	OutputDir   string

	// Framework & Logger
	Framework Framework
	Logger    Logger

	// Optional packages
	Viper bool
	Redis bool
	Kafka bool
	NATS  bool

	// DevOps
	Docker   bool
	Makefile bool
	GitHub   bool
	Lint     bool
	Swagger  bool
}

// FrameworkImport returns the Go module import path for the selected framework.
func (c *ProjectConfig) FrameworkImport() string {
	switch c.Framework {
	case FrameworkGin:
		return "github.com/gin-gonic/gin"
	case FrameworkFiber:
		return "github.com/gofiber/fiber/v2"
	case FrameworkChi:
		return "github.com/go-chi/chi/v5"
	case FrameworkEcho:
		return "github.com/labstack/echo/v4"
	case FrameworkGorilla:
		return "github.com/gorilla/mux"
	default:
		return ""
	}
}

// LoggerImport returns the Go module import path for the selected logger.
func (c *ProjectConfig) LoggerImport() string {
	switch c.Logger {
	case LoggerZerolog:
		return "github.com/rs/zerolog"
	case LoggerZap:
		return "go.uber.org/zap"
	default:
		return "" // slog is stdlib
	}
}

// HasInfra returns true if any infrastructure package is selected.
func (c *ProjectConfig) HasInfra() bool {
	return c.Redis || c.Kafka || c.NATS
}

// Validate validates the project configuration.
func (c *ProjectConfig) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if c.ModuleName == "" {
		return fmt.Errorf("module name is required")
	}
	return nil
}
