package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/iyashjayesh/goscaf/internal/config"
)

// Run runs the interactive prompt flow and returns a populated ProjectConfig.
func Run(projectName string) (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{
		ProjectName: projectName,
	}

	// 1. Module name
	moduleDefault := fmt.Sprintf("github.com/your-org/%s", projectName)
	if err := survey.AskOne(&survey.Input{
		Message: "Module name:",
		Default: moduleDefault,
	}, &cfg.ModuleName, survey.WithValidator(survey.Required)); err != nil {
		return nil, fmt.Errorf("ask module name: %w", err)
	}

	// 2. Go version
	goVersionStr := "1.25.0"
	if err := survey.AskOne(&survey.Select{
		Message: "Go version:",
		Options: []string{"1.25.0", "1.24.0", "1.23"},
		Default: "1.25.0",
	}, &goVersionStr); err != nil {
		return nil, fmt.Errorf("ask go version: %w", err)
	}
	cfg.GoVersion = goVersionStr

	// 3. HTTP framework
	frameworkStr := "gin"
	if err := survey.AskOne(&survey.Select{
		Message: "HTTP framework:",
		Options: []string{"gin", "fiber", "chi", "echo", "gorilla/mux", "none"},
		Default: "gin",
	}, &frameworkStr); err != nil {
		return nil, fmt.Errorf("ask framework: %w", err)
	}
	if frameworkStr == "gorilla/mux" {
		frameworkStr = "gorilla"
	}
	cfg.Framework = config.Framework(frameworkStr)

	// 4. Structured logger
	loggerStr := "slog"
	if err := survey.AskOne(&survey.Select{
		Message: "Structured logger:",
		Options: []string{"slog (stdlib)", "zerolog", "zap"},
		Default: "slog (stdlib)",
	}, &loggerStr); err != nil {
		return nil, fmt.Errorf("ask logger: %w", err)
	}
	switch loggerStr {
	case "slog (stdlib)":
		cfg.Logger = config.LoggerSlog
	case "zerolog":
		cfg.Logger = config.LoggerZerolog
	case "zap":
		cfg.Logger = config.LoggerZap
	}

	// 5. Viper
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Viper for config & env management?",
		Default: true,
	}, &cfg.Viper); err != nil {
		return nil, fmt.Errorf("ask viper: %w", err)
	}

	// 6. Redis
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Redis client (go-redis)?",
		Default: false,
	}, &cfg.Redis); err != nil {
		return nil, fmt.Errorf("ask redis: %w", err)
	}

	// 7. Kafka
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Kafka client (franz-go)?",
		Default: false,
	}, &cfg.Kafka); err != nil {
		return nil, fmt.Errorf("ask kafka: %w", err)
	}

	// 8. NATS
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add NATS client?",
		Default: false,
	}, &cfg.NATS); err != nil {
		return nil, fmt.Errorf("ask nats: %w", err)
	}

	// 9. Database driver
	dbStr := "none"
	if err := survey.AskOne(&survey.Select{
		Message: "Database driver:",
		Options: []string{"none", "postgres", "mysql", "sqlite", "mongo", "gorm"},
		Default: "none",
	}, &dbStr); err != nil {
		return nil, fmt.Errorf("ask database: %w", err)
	}
	cfg.Database = config.Database(dbStr)

	// 10. Dockerfile + docker-compose
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Dockerfile + docker-compose?",
		Default: true,
	}, &cfg.Docker); err != nil {
		return nil, fmt.Errorf("ask docker: %w", err)
	}

	// 10. Makefile
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Makefile?",
		Default: true,
	}, &cfg.Makefile); err != nil {
		return nil, fmt.Errorf("ask makefile: %w", err)
	}

	// 11. GitHub Actions CI
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add GitHub Actions CI?",
		Default: true,
	}, &cfg.GitHub); err != nil {
		return nil, fmt.Errorf("ask github: %w", err)
	}

	// 12. golangci-lint config
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add golangci-lint config?",
		Default: true,
	}, &cfg.Lint); err != nil {
		return nil, fmt.Errorf("ask lint: %w", err)
	}

	// 13. Swagger/OpenAPI scaffold
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Swagger/OpenAPI scaffold?",
		Default: false,
	}, &cfg.Swagger); err != nil {
		return nil, fmt.Errorf("ask swagger: %w", err)
	}

	return cfg, nil
}
