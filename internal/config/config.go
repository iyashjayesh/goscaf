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

// Database represents the database driver choice.
type Database string

const (
	DBPostgres Database = "postgres"
	DBMySQL    Database = "mysql"
	DBSQLite   Database = "sqlite"
	DBMongo    Database = "mongo"
	DBGORM     Database = "gorm"
	DBNone     Database = "none"
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
	Viper    bool
	Redis    bool
	Kafka    bool
	NATS     bool
	Database Database

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

// DBImport returns the Go module import path for the selected database driver.
func (c *ProjectConfig) DBImport() string {
	switch c.Database {
	case DBPostgres:
		return "github.com/jackc/pgx/v5"
	case DBMySQL:
		return "github.com/go-sql-driver/mysql"
	case DBSQLite:
		return "modernc.org/sqlite"
	case DBMongo:
		return "go.mongodb.org/mongo-driver/mongo"
	case DBGORM:
		return "gorm.io/gorm"
	default:
		return ""
	}
}

// HasDB returns true if a real database driver is selected.
func (c *ProjectConfig) HasDB() bool {
	return c.Database != "" && c.Database != DBNone
}

// HasInfra returns true if any infrastructure package is selected.
func (c *ProjectConfig) HasInfra() bool {
	return c.Redis || c.Kafka || c.NATS || c.HasDB()
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
