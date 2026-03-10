package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iyashjayesh/gostart/internal/config"
	"github.com/iyashjayesh/gostart/internal/generator"
	"github.com/iyashjayesh/gostart/internal/prompt"
)

var (
	flagModule    string
	flagGoVersion string
	flagFramework string
	flagLogger    string
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

		var cfg *config.ProjectConfig

		if flagDefaults {
			// Use recommended defaults without prompting
			cfg = &config.ProjectConfig{
				ProjectName: projectName,
				ModuleName:  fmt.Sprintf("github.com/your-org/%s", projectName),
				GoVersion:   "1.23",
				Framework:   config.FrameworkGin,
				Logger:      config.LoggerSlog,
				Viper:       true,
				Redis:       false,
				Kafka:       false,
				NATS:        false,
				Docker:      true,
				Makefile:    true,
				GitHub:      true,
				Lint:        true,
				Swagger:     false,
			}
		} else if cmd.Flags().Changed("framework") || cmd.Flags().Changed("module") ||
			cmd.Flags().Changed("go-version") || cmd.Flags().Changed("logger") {
			// Flags provided - use flag-driven mode (merge with defaults)
			cfg = &config.ProjectConfig{
				ProjectName: projectName,
				ModuleName:  flagModule,
				GoVersion:   flagGoVersion,
				Framework:   config.Framework(flagFramework),
				Logger:      config.Logger(flagLogger),
				Viper:       flagViper,
				Redis:       flagRedis,
				Kafka:       flagKafka,
				NATS:        flagNATS,
				Docker:      flagDocker,
				Makefile:    flagMakefile,
				GitHub:      flagGitHub,
				Lint:        flagLint,
				Swagger:     flagSwagger,
			}
			if cfg.ModuleName == "" {
				cfg.ModuleName = fmt.Sprintf("github.com/your-org/%s", projectName)
			}
		} else {
			// Interactive mode
			var err error
			cfg, err = prompt.Run(projectName)
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
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

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&flagModule, "module", "", "Go module path")
	initCmd.Flags().StringVar(&flagGoVersion, "go-version", "1.23", "Go version (1.21, 1.22, 1.23)")
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
	initCmd.Flags().BoolVar(&flagDefaults, "defaults", false, "Skip all prompts, use recommended defaults")
	initCmd.Flags().StringVar(&flagOutput, "output", "", "Output directory (default: current dir)")
}
