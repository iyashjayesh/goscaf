package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iyashjayesh/goscaf/internal/config"
	"github.com/iyashjayesh/goscaf/internal/generator"
	"github.com/iyashjayesh/goscaf/internal/prompt"
	"github.com/iyashjayesh/goscaf/internal/userconfig"
)

var (
	flagModule    string
	flagGoVersion string
	flagFramework string
	flagLogger    string
	flagDB        string
	flagViper     bool
	flagRedis     bool
	flagKafka     bool
	flagNATS      bool
	flagDocker    bool
	flagMakefile  bool
	flagGitHub    bool
	flagLint      bool
	flagSwagger   bool
	flagDefaults  bool
	flagOutput    string
	flagGitRepo   bool
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Scaffold a new enterprise-grade Go project",
	Long: `Scaffold a new enterprise-grade Go project with your choice of:
  • HTTP framework (gin, fiber, chi, echo, gorilla/mux)
  • Structured logger (slog, zerolog, zap)
  • Infrastructure clients (Redis, Kafka, NATS)
  • DevOps tooling (Docker, Makefile, GitHub Actions, golangci-lint)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Load .goscaf.yaml (global then local, local wins). A missing file is
		// not an error — uc will be nil and prompts fall back to hardcoded defaults.
		uc, err := userconfig.Load()
		if err != nil {
			color.Yellow("  ⚠ could not read .goscaf.yaml: %v", err)
		}

		var cfg *config.ProjectConfig

		if flagDefaults {
			// Use recommended defaults without prompting
			cfg = &config.ProjectConfig{
				ProjectName: projectName,
				ModuleName:  fmt.Sprintf("github.com/your-org/%s", projectName),
				GoVersion:   "1.25.0",
				Framework:   config.FrameworkGin,
				Logger:      config.LoggerSlog,
				Database:    config.Database(flagDB),
				Viper:       true,
				Redis:       false,
				Kafka:       false,
				NATS:        false,
				Docker:      true,
				Makefile:    true,
				GitHub:      true,
				Lint:        true,
				Swagger:     false,
				GitRepo:     false,
			}
		} else if cmd.Flags().Changed("framework") || cmd.Flags().Changed("module") ||
			cmd.Flags().Changed("go-version") || cmd.Flags().Changed("logger") ||
			cmd.Flags().Changed("db") {
			// Flags provided — start from flag values, then fill in any field the
			// user did NOT explicitly set from .goscaf.yaml (flags always win).
			cfg = &config.ProjectConfig{
				ProjectName: projectName,
				ModuleName:  flagModule,
				GoVersion:   flagGoVersion,
				Framework:   config.Framework(flagFramework),
				Logger:      config.Logger(flagLogger),
				Database:    config.Database(flagDB),
				Viper:       flagViper,
				Redis:       flagRedis,
				Kafka:       flagKafka,
				NATS:        flagNATS,
				Docker:      flagDocker,
				Makefile:    flagMakefile,
				GitHub:      flagGitHub,
				Lint:        flagLint,
				Swagger:     flagSwagger,
				GitRepo:     flagGitRepo,
			}
			if cfg.ModuleName == "" {
				if uc != nil && uc.ModulePrefix != "" {
					cfg.ModuleName = uc.ModulePrefix + "/" + projectName
				} else {
					cfg.ModuleName = fmt.Sprintf("github.com/your-org/%s", projectName)
				}
			}
			if uc != nil {
				applyUserConfig(cmd, cfg, uc)
			}
		} else {
			// Interactive mode — userconfig pre-fills prompt defaults.
			var promptErr error
			cfg, promptErr = prompt.Run(projectName, uc)
			if promptErr != nil {
				return fmt.Errorf("prompt failed: %w", promptErr)
			}
		}

		// Set output directory
		outputDir := flagOutput
		if outputDir == "" {
			outputDir = "."
		}
		cfg.OutputDir = filepath.Join(outputDir, projectName)

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		// Print config summary
		fmt.Println()
		color.HiCyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		color.HiWhite("  Scaffolding project: %s", color.HiYellowString(cfg.ProjectName))
		color.HiWhite("  Module:    %s", cfg.ModuleName)
		color.HiWhite("  Go:        %s", cfg.GoVersion)
		color.HiWhite("  Framework: %s", string(cfg.Framework))
		color.HiWhite("  Logger:    %s", string(cfg.Logger))
		color.HiWhite("  Database:  %s", string(cfg.Database))
		color.HiCyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		// Run generator
		gen := generator.New(cfg)
		if err := gen.Run(); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		// Print next steps
		fmt.Println()
		color.HiGreen("  ✔ Project created successfully!")
		fmt.Println()
		color.HiCyan("  Next steps:")
		color.HiWhite("    cd %s", cfg.ProjectName)
		color.HiWhite("    cp .env.example .env")
		if cfg.Docker {
			color.HiWhite("    make docker-up")
		}
		color.HiWhite("    make run")
		fmt.Println()

		os.Exit(0)
		return nil
	},
}

// applyUserConfig fills cfg fields from uc for any flag the user did not
// explicitly pass on the command line. CLI flags always take precedence.
func applyUserConfig(cmd *cobra.Command, cfg *config.ProjectConfig, uc *userconfig.UserConfig) {
	if !cmd.Flags().Changed("go-version") && uc.GoVersion != "" {
		cfg.GoVersion = uc.GoVersion
	}
	if !cmd.Flags().Changed("framework") && uc.Framework != "" {
		cfg.Framework = config.Framework(uc.Framework)
	}
	if !cmd.Flags().Changed("logger") && uc.Logger != "" {
		cfg.Logger = config.Logger(uc.Logger)
	}
	if !cmd.Flags().Changed("db") && uc.DB != "" {
		cfg.Database = config.Database(uc.DB)
	}
	if !cmd.Flags().Changed("viper") && uc.Viper != nil {
		cfg.Viper = *uc.Viper
	}
	if !cmd.Flags().Changed("redis") && uc.Redis != nil {
		cfg.Redis = *uc.Redis
	}
	if !cmd.Flags().Changed("kafka") && uc.Kafka != nil {
		cfg.Kafka = *uc.Kafka
	}
	if !cmd.Flags().Changed("nats") && uc.NATS != nil {
		cfg.NATS = *uc.NATS
	}
	if !cmd.Flags().Changed("docker") && uc.Docker != nil {
		cfg.Docker = *uc.Docker
	}
	if !cmd.Flags().Changed("makefile") && uc.Makefile != nil {
		cfg.Makefile = *uc.Makefile
	}
	if !cmd.Flags().Changed("github") && uc.GitHub != nil {
		cfg.GitHub = *uc.GitHub
	}
	if !cmd.Flags().Changed("lint") && uc.Lint != nil {
		cfg.Lint = *uc.Lint
	}
	if !cmd.Flags().Changed("swagger") && uc.Swagger != nil {
		cfg.Swagger = *uc.Swagger
	}
	if !cmd.Flags().Changed("git-repo") && uc.GitRepo != nil {
		cfg.GitRepo = *uc.GitRepo
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&flagModule, "module", "", "Go module path")
	initCmd.Flags().StringVar(&flagGoVersion, "go-version", "1.25.0", "Go version (1.23, 1.24, 1.25)")
	initCmd.Flags().StringVar(&flagFramework, "framework", "gin", "HTTP framework (gin|fiber|chi|echo|gorilla|none)")
	initCmd.Flags().StringVar(&flagLogger, "logger", "slog", "Structured logger (slog|zerolog|zap)")
	initCmd.Flags().BoolVar(&flagViper, "viper", true, "Add Viper for config & env management")
	initCmd.Flags().BoolVar(&flagRedis, "redis", false, "Add Redis client (go-redis)")
	initCmd.Flags().BoolVar(&flagKafka, "kafka", false, "Add Kafka client (franz-go)")
	initCmd.Flags().BoolVar(&flagNATS, "nats", false, "Add NATS client")
	initCmd.Flags().BoolVar(&flagDocker, "docker", true, "Add Dockerfile + docker-compose")
	initCmd.Flags().BoolVar(&flagMakefile, "makefile", true, "Add Makefile")
	initCmd.Flags().BoolVar(&flagGitHub, "github", true, "Add GitHub Actions CI")
	initCmd.Flags().BoolVar(&flagLint, "lint", true, "Add golangci-lint config")
	initCmd.Flags().BoolVar(&flagSwagger, "swagger", false, "Add Swagger/OpenAPI scaffold")
	initCmd.Flags().BoolVar(&flagGitRepo, "git-repo", false, "Initialize a new git repository for this project")
	initCmd.Flags().StringVar(&flagDB, "db", "none", "Database driver (postgres|mysql|sqlite|mongo|gorm|none)")
	initCmd.Flags().BoolVar(&flagDefaults, "defaults", false, "Skip all prompts, use recommended defaults")
	initCmd.Flags().StringVar(&flagOutput, "output", "", "Output directory (default: current dir)")

}
