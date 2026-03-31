package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type ServiceGenerator struct {
	projectDir string
	moduleName string
	framework  string
}

type ServiceInfo struct {
	InputName     string
	DirectoryName string
	FileBaseName  string
	PackageName   string
	StructName    string
	ModuleName    string
}

func NewServiceGenerator(projectDir string) (*ServiceGenerator, error) {
	moduleName, err := readModuleName(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return nil, err
	}

	framework, err := detectFramework(filepath.Join(projectDir, "internal/server/server.go"))
	if err != nil {
		return nil, err
	}

	return &ServiceGenerator{
		projectDir: projectDir,
		moduleName: moduleName,
		framework:  framework,
	}, nil
}

func (g *ServiceGenerator) Run(rawName string) (*ServiceInfo, error) {
	info, err := newServiceInfo(rawName, g.moduleName)
	if err != nil {
		return nil, err
	}

	steps := []struct {
		path string
		tmpl string
	}{
		{
			path: filepath.Join("internal", info.DirectoryName, "service.go"),
			tmpl: serviceTemplate,
		},
		{
			path: filepath.Join("internal", "handler", info.FileBaseName+"_handler.go"),
			tmpl: handlerTemplate(g.framework),
		},
	}

	for _, step := range steps {
		fullPath := filepath.Join(g.projectDir, step.path)
		if _, err := os.Stat(fullPath); err == nil {
			return nil, fmt.Errorf("%s already exists", step.path)
		}
		if err := writeTemplateFile(fullPath, step.tmpl, info); err != nil {
			return nil, fmt.Errorf("generate %s: %w", step.path, err)
		}
		color.HiGreen("  ✓ %s", step.path)
	}

	return info, nil
}

func newServiceInfo(rawName, moduleName string) (*ServiceInfo, error) {
	fileBase := sanitizeServiceName(rawName)
	if fileBase == "" {
		return nil, fmt.Errorf("service name %q is invalid", rawName)
	}

	structName := toCamel(fileBase)
	return &ServiceInfo{
		InputName:     rawName,
		DirectoryName: fileBase,
		FileBaseName:  fileBase,
		PackageName:   fileBase,
		StructName:    structName,
		ModuleName:    moduleName,
	}, nil
}

func sanitizeServiceName(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(normalized, "_")
	normalized = regexp.MustCompile(`_+`).ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return ""
	}
	if normalized[0] >= '0' && normalized[0] <= '9' {
		normalized = "service_" + normalized
	}
	return normalized
}

func toCamel(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(part[1:])
		}
	}
	return b.String()
}

func readModuleName(goModPath string) (string, error) {
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}

	matches := regexp.MustCompile(`(?m)^module\s+(\S+)\s*$`).FindSubmatch(content)
	if len(matches) != 2 {
		return "", fmt.Errorf("could not determine module name from go.mod")
	}

	return string(matches[1]), nil
}

func detectFramework(serverPath string) (string, error) {
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return "", fmt.Errorf("read internal/server/server.go: %w", err)
	}

	serverSource := string(content)
	switch {
	case strings.Contains(serverSource, `"github.com/gin-gonic/gin"`):
		return "gin", nil
	case strings.Contains(serverSource, `"github.com/gofiber/fiber/v2"`):
		return "fiber", nil
	case strings.Contains(serverSource, `"github.com/go-chi/chi/v5"`):
		return "chi", nil
	case strings.Contains(serverSource, `"github.com/labstack/echo/v4"`):
		return "echo", nil
	case strings.Contains(serverSource, `"github.com/gorilla/mux"`):
		return "gorilla", nil
	default:
		return "none", nil
	}
}

func writeTemplateFile(fullPath, tmplStr string, data any) error {
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(fullPath), err)
	}

	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if err := os.WriteFile(fullPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

const serviceTemplate = `package {{.PackageName}}

import "context"

// Service holds the business logic for {{.InputName}}.
type Service struct{}

// NewService creates a new {{.StructName}} service.
func NewService() *Service {
	return &Service{}
}

// Health returns a simple readiness message for the service.
func (s *Service) Health(ctx context.Context) (string, error) {
	_ = ctx
	return "{{.InputName}} service is healthy", nil
}
`

func handlerTemplate(framework string) string {
	switch framework {
	case "gin":
		return handlerGinTemplate
	case "fiber":
		return handlerFiberTemplate
	case "chi":
		return handlerChiTemplate
	case "echo":
		return handlerEchoTemplate
	case "gorilla":
		return handlerGorillaTemplate
	case "none":
		return handlerHTTPTemplate
	default:
		return handlerHTTPTemplate
	}
}

const handlerGinTemplate = `package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided router group.
func (h *{{.StructName}}Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/{{.DirectoryName}}/health", h.Health)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(c *gin.Context) {
	message, err := h.service.Health(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
`

const handlerFiberTemplate = `package handler

import (
	"github.com/gofiber/fiber/v2"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided router.
func (h *{{.StructName}}Handler) Register(rg fiber.Router) {
	rg.Get("/{{.DirectoryName}}/health", h.Health)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(c *fiber.Ctx) error {
	message, err := h.service.Health(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
}
`

const handlerEchoTemplate = `package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided group.
func (h *{{.StructName}}Handler) Register(rg *echo.Group) {
	rg.GET("/{{.DirectoryName}}/health", h.Health)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(c echo.Context) error {
	message, err := h.service.Health(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}
`

const handlerHTTPTemplate = `package handler

import (
	"encoding/json"
	"net/http"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided mux.
func (h *{{.StructName}}Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/{{.DirectoryName}}/health", h.Health)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(w http.ResponseWriter, r *http.Request) {
	message, err := h.service.Health(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
`

const handlerChiTemplate = `package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided router.
func (h *{{.StructName}}Handler) Register(r chi.Router) {
	r.Get("/{{.DirectoryName}}/health", h.Health)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(w http.ResponseWriter, r *http.Request) {
	message, err := h.service.Health(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
`

const handlerGorillaTemplate = `package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"{{.ModuleName}}/internal/{{.DirectoryName}}"
)

// {{.StructName}}Handler exposes the HTTP handlers for the {{.InputName}} domain.
type {{.StructName}}Handler struct {
	service *{{.PackageName}}.Service
}

// New{{.StructName}}Handler creates a new {{.StructName}}Handler.
func New{{.StructName}}Handler(service *{{.PackageName}}.Service) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{service: service}
}

// Register mounts {{.InputName}} routes under the provided router.
func (h *{{.StructName}}Handler) Register(r *mux.Router) {
	r.HandleFunc("/{{.DirectoryName}}/health", h.Health).Methods(http.MethodGet)
}

// Health responds with the {{.InputName}} service health status.
func (h *{{.StructName}}Handler) Health(w http.ResponseWriter, r *http.Request) {
	message, err := h.service.Health(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
`
