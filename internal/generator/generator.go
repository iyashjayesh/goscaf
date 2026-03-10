package generator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/fatih/color"

	"github.com/iyashjayesh/gostart/internal/config"
	"github.com/iyashjayesh/gostart/internal/templates"
)

// Generator orchestrates all file writes for a project.
type Generator struct {
	cfg *config.ProjectConfig
}

// New creates a new Generator.
func New(cfg *config.ProjectConfig) *Generator {
	return &Generator{cfg: cfg}
}

// Run executes all generation steps.
func (g *Generator) Run() error {
	steps := g.buildSteps()
	for _, step := range steps {
		if step.skip {
			continue
		}
		if err := g.writeFile(step.path, step.tmpl); err != nil {
			color.Red("  ✗ %s - %v", step.path, err)
			return err
		}
		color.HiGreen("  ✓ %s", step.path)
	}

	// Run go mod tidy
	fmt.Println()
	color.HiBlue("  → Running go mod tidy...")
	tidyCmd := exec.CommandContext(context.Background(), "go", "mod", "tidy")
	tidyCmd.Dir = g.cfg.OutputDir
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		color.Yellow("  ⚠ go mod tidy failed (you may need to run it manually): %v", err)
	} else {
		color.HiGreen("  ✓ go mod tidy")
	}

	return nil
}

type step struct {
	path string
	tmpl string
	skip bool
}

func (g *Generator) buildSteps() []step {
	cfg := g.cfg
	steps := []step{
		// Base files
		{path: "cmd/main.go", tmpl: templates.MainGo},
		{path: "go.mod", tmpl: templates.GoMod},
		{path: ".gitignore", tmpl: templates.GitIgnore},
		{path: ".env.example", tmpl: templates.EnvExample},
		{path: "internal/config/config.go", tmpl: templates.ConfigGo(cfg.Viper)},
		// Handler
		{path: "internal/handler/handler.go", tmpl: templates.HandlerGo},
		// Server
		{path: "internal/server/server.go", tmpl: templates.ServerGo(cfg.Framework)},
	}

	// Infrastructure packages
	if cfg.Redis {
		steps = append(steps, step{path: "pkg/redis/redis.go", tmpl: templates.RedisGo})
	}
	if cfg.Kafka {
		steps = append(steps, step{path: "pkg/kafka/kafka.go", tmpl: templates.KafkaGo})
	}
	if cfg.NATS {
		steps = append(steps, step{path: "pkg/nats/nats.go", tmpl: templates.NatsGo})
	}

	// DevOps files
	if cfg.Docker {
		steps = append(steps, step{path: "Dockerfile", tmpl: templates.DockerfileTemplate})
		steps = append(steps, step{path: "docker-compose.yml", tmpl: templates.DockerCompose})
	}
	if cfg.Makefile {
		steps = append(steps, step{path: "Makefile", tmpl: templates.MakefileTemplate})
	}
	if cfg.GitHub {
		steps = append(steps, step{path: ".github/workflows/ci.yml", tmpl: templates.GitHubCI})
	}
	if cfg.Lint {
		steps = append(steps, step{path: ".golangci.yml", tmpl: templates.GolangCI})
	}
	if cfg.Swagger {
		steps = append(steps, step{path: "docs/swagger.yaml", tmpl: templates.SwaggerYAML})
	}

	return steps
}

// writeFile renders the template and writes to the target path.
func (g *Generator) writeFile(relPath, tmplStr string) error {
	fullPath := filepath.Join(g.cfg.OutputDir, relPath)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"eq": func(a, b interface{}) bool { return fmt.Sprint(a) == fmt.Sprint(b) },
	}).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.cfg); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if err := os.WriteFile(fullPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
