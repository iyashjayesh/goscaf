package userconfig

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// UserConfig holds optional defaults loaded from .goscaf.yaml.
// String fields use empty string as "not set".
// Bool fields use *bool so nil (absent) is distinct from false (explicitly disabled).
type UserConfig struct {
	ModulePrefix string `yaml:"module_prefix"`
	GoVersion    string `yaml:"go_version"`
	Framework    string `yaml:"framework"`
	Logger       string `yaml:"logger"`
	DB           string `yaml:"db"`
	Viper        *bool  `yaml:"viper"`
	Redis        *bool  `yaml:"redis"`
	Kafka        *bool  `yaml:"kafka"`
	NATS         *bool  `yaml:"nats"`
	Docker       *bool  `yaml:"docker"`
	Makefile     *bool  `yaml:"makefile"`
	GitHub       *bool  `yaml:"github"`
	Lint         *bool  `yaml:"lint"`
	Swagger      *bool  `yaml:"swagger"`
	GitRepo      *bool  `yaml:"git_repo"`
}

// Load reads and merges ~/.goscaf.yaml (global) and ./.goscaf.yaml (local).
// Local values take precedence. Returns nil if neither file exists.
func Load() (*UserConfig, error) {
	home, _ := os.UserHomeDir()
	globalPath := ""
	if home != "" {
		globalPath = filepath.Join(home, ".goscaf.yaml")
	}
	return loadFrom(globalPath, ".goscaf.yaml")
}

// loadFrom is the testable core: it reads from explicit paths rather than
// deriving them from the environment.
func loadFrom(globalPath, localPath string) (*UserConfig, error) {
	global, err := readFile(globalPath)
	if err != nil {
		return nil, err
	}

	local, err := readFile(localPath)
	if err != nil {
		return nil, err
	}

	if global == nil && local == nil {
		return nil, nil
	}

	return merge(global, local), nil
}

// readFile parses a single .goscaf.yaml file. Returns nil (not an error) when
// the file does not exist, so callers can treat absence as "no config".
func readFile(path string) (*UserConfig, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg UserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// merge combines global and local, with local non-zero values taking precedence.
func merge(global, local *UserConfig) *UserConfig {
	result := &UserConfig{}

	if global != nil {
		*result = *global
	}

	if local == nil {
		return result
	}

	if local.ModulePrefix != "" {
		result.ModulePrefix = local.ModulePrefix
	}
	if local.GoVersion != "" {
		result.GoVersion = local.GoVersion
	}
	if local.Framework != "" {
		result.Framework = local.Framework
	}
	if local.Logger != "" {
		result.Logger = local.Logger
	}
	if local.DB != "" {
		result.DB = local.DB
	}
	if local.Viper != nil {
		result.Viper = local.Viper
	}
	if local.Redis != nil {
		result.Redis = local.Redis
	}
	if local.Kafka != nil {
		result.Kafka = local.Kafka
	}
	if local.NATS != nil {
		result.NATS = local.NATS
	}
	if local.Docker != nil {
		result.Docker = local.Docker
	}
	if local.Makefile != nil {
		result.Makefile = local.Makefile
	}
	if local.GitHub != nil {
		result.GitHub = local.GitHub
	}
	if local.Lint != nil {
		result.Lint = local.Lint
	}
	if local.Swagger != nil {
		result.Swagger = local.Swagger
	}
	if local.GitRepo != nil {
		result.GitRepo = local.GitRepo
	}

	return result
}

// BoolVal returns the dereferenced value of b, or fallback if b is nil.
func BoolVal(b *bool, fallback bool) bool {
	if b == nil {
		return fallback
	}
	return *b
}
