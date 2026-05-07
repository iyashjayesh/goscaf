package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/iyashjayesh/goscaf/internal/config"
	"github.com/iyashjayesh/goscaf/internal/userconfig"
)

// Run runs the interactive prompt flow and returns a populated ProjectConfig.
// If uc is non-nil its values are used as pre-filled defaults; the user can
// still change any answer at the prompt.
func Run(projectName string, uc *userconfig.UserConfig) (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{
		ProjectName: projectName,
	}

	// 1. Module name
	moduleDefault := fmt.Sprintf("github.com/your-org/%s", projectName)
	if uc != nil && uc.ModulePrefix != "" {
		moduleDefault = uc.ModulePrefix + "/" + projectName
	}
	if err := survey.AskOne(&survey.Input{
		Message: "Module name:",
		Default: moduleDefault,
	}, &cfg.ModuleName, survey.WithValidator(survey.Required)); err != nil {
		return nil, fmt.Errorf("ask module name: %w", err)
	}

	// 2. Go version
	goVersionDefault := "1.25.0"
	if uc != nil && uc.GoVersion != "" {
		goVersionDefault = uc.GoVersion
	}
	goVersionStr := goVersionDefault
	if err := survey.AskOne(&survey.Select{
		Message: "Go version:",
		Options: []string{"1.25.0", "1.24.0", "1.23"},
		Default: goVersionDefault,
	}, &goVersionStr); err != nil {
		return nil, fmt.Errorf("ask go version: %w", err)
	}
	cfg.GoVersion = goVersionStr

	// 3. HTTP framework
	// The select option for gorilla is "gorilla/mux" but the config value is "gorilla".
	frameworkDefault := "gin"
	if uc != nil && uc.Framework != "" {
		f := uc.Framework
		if f == "gorilla" {
			f = "gorilla/mux"
		}
		frameworkDefault = f
	}
	frameworkStr := frameworkDefault
	if err := survey.AskOne(&survey.Select{
		Message: "HTTP framework:",
		Options: []string{"gin", "fiber", "chi", "echo", "gorilla/mux", "none"},
		Default: frameworkDefault,
	}, &frameworkStr); err != nil {
		return nil, fmt.Errorf("ask framework: %w", err)
	}
	if frameworkStr == "gorilla/mux" {
		frameworkStr = "gorilla"
	}
	cfg.Framework = config.Framework(frameworkStr)

	// 4. Structured logger
	// The select option for slog is "slog (stdlib)" but the config value is "slog".
	loggerDefault := "slog (stdlib)"
	if uc != nil && uc.Logger != "" {
		switch uc.Logger {
		case "slog":
			loggerDefault = "slog (stdlib)"
		default:
			loggerDefault = uc.Logger
		}
	}
	loggerStr := loggerDefault
	if err := survey.AskOne(&survey.Select{
		Message: "Structured logger:",
		Options: []string{"slog (stdlib)", "zerolog", "zap"},
		Default: loggerDefault,
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
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Viper }), true),
	}, &cfg.Viper); err != nil {
		return nil, fmt.Errorf("ask viper: %w", err)
	}

	// 6. Redis
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Redis client (go-redis)?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Redis }), false),
	}, &cfg.Redis); err != nil {
		return nil, fmt.Errorf("ask redis: %w", err)
	}

	// 7. Kafka
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Kafka client (franz-go)?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Kafka }), false),
	}, &cfg.Kafka); err != nil {
		return nil, fmt.Errorf("ask kafka: %w", err)
	}

	// 8. NATS
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add NATS client?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.NATS }), false),
	}, &cfg.NATS); err != nil {
		return nil, fmt.Errorf("ask nats: %w", err)
	}

	// 9. Database driver
	dbDefault := "none"
	if uc != nil && uc.DB != "" {
		dbDefault = uc.DB
	}
	dbStr := dbDefault
	if err := survey.AskOne(&survey.Select{
		Message: "Database driver:",
		Options: []string{"none", "postgres", "mysql", "sqlite", "mongo", "gorm"},
		Default: dbDefault,
	}, &dbStr); err != nil {
		return nil, fmt.Errorf("ask database: %w", err)
	}
	cfg.Database = config.Database(dbStr)

	// 10. Dockerfile + docker-compose
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Dockerfile + docker-compose?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Docker }), true),
	}, &cfg.Docker); err != nil {
		return nil, fmt.Errorf("ask docker: %w", err)
	}

	// 11. Makefile
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Makefile?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Makefile }), true),
	}, &cfg.Makefile); err != nil {
		return nil, fmt.Errorf("ask makefile: %w", err)
	}

	// 12. GitHub Actions CI
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add GitHub Actions CI?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.GitHub }), true),
	}, &cfg.GitHub); err != nil {
		return nil, fmt.Errorf("ask github: %w", err)
	}

	// 13. golangci-lint config
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add golangci-lint config?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Lint }), true),
	}, &cfg.Lint); err != nil {
		return nil, fmt.Errorf("ask lint: %w", err)
	}

	// 14. Swagger/OpenAPI scaffold
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add Swagger/OpenAPI scaffold?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.Swagger }), false),
	}, &cfg.Swagger); err != nil {
		return nil, fmt.Errorf("ask swagger: %w", err)
	}

	// 15. Git repository
	if err := survey.AskOne(&survey.Confirm{
		Message: "Initialize github repository?",
		Default: userconfig.BoolVal(ucBool(uc, func(u *userconfig.UserConfig) *bool { return u.GitRepo }), false),
	}, &cfg.GitRepo); err != nil {
		return nil, fmt.Errorf("failed to initialize git repository : %w", err)
	}

	return cfg, nil
}

// ucBool safely extracts a *bool field from uc, returning nil when uc is nil.
func ucBool(uc *userconfig.UserConfig, fn func(*userconfig.UserConfig) *bool) *bool {
	if uc == nil {
		return nil
	}
	return fn(uc)
}
