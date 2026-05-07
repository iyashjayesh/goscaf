package userconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// boolPtr is a test helper to get a pointer to a bool literal.
func boolPtr(v bool) *bool { return &v }

// writeYAML writes content to a temp file and returns its path.
func writeYAML(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeYAML: %v", err)
	}
	return path
}

func TestLoadFromNoFiles(t *testing.T) {
	dir := t.TempDir()
	got, err := loadFrom(
		filepath.Join(dir, "global.yaml"),
		filepath.Join(dir, "local.yaml"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil when no files exist, got %+v", got)
	}
}

func TestLoadFromGlobalOnly(t *testing.T) {
	dir := t.TempDir()
	global := writeYAML(t, dir, "global.yaml", `
module_prefix: github.com/global-org
framework: gin
logger: zap
viper: true
`)

	got, err := loadFrom(global, filepath.Join(dir, "local.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil config")
	}
	if got.ModulePrefix != "github.com/global-org" {
		t.Errorf("ModulePrefix: got %q, want %q", got.ModulePrefix, "github.com/global-org")
	}
	if got.Framework != "gin" {
		t.Errorf("Framework: got %q, want %q", got.Framework, "gin")
	}
	if got.Logger != "zap" {
		t.Errorf("Logger: got %q, want %q", got.Logger, "zap")
	}
	if got.Viper == nil || !*got.Viper {
		t.Errorf("Viper: expected true, got %v", got.Viper)
	}
}

func TestLoadFromLocalOnly(t *testing.T) {
	dir := t.TempDir()
	local := writeYAML(t, dir, "local.yaml", `
framework: fiber
docker: false
`)

	got, err := loadFrom(filepath.Join(dir, "global.yaml"), local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil config")
	}
	if got.Framework != "fiber" {
		t.Errorf("Framework: got %q, want %q", got.Framework, "fiber")
	}
	if got.Docker == nil || *got.Docker {
		t.Errorf("Docker: expected false pointer, got %v", got.Docker)
	}
}

func TestLocalOverridesGlobalString(t *testing.T) {
	dir := t.TempDir()
	global := writeYAML(t, dir, "global.yaml", `
framework: gin
logger: slog
module_prefix: github.com/global-org
`)
	local := writeYAML(t, dir, "local.yaml", `
framework: fiber
module_prefix: github.com/local-org
`)

	got, err := loadFrom(global, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Framework != "fiber" {
		t.Errorf("Framework: got %q, want fiber (local should win)", got.Framework)
	}
	if got.ModulePrefix != "github.com/local-org" {
		t.Errorf("ModulePrefix: got %q, want github.com/local-org", got.ModulePrefix)
	}
	// logger not set in local — global value should survive
	if got.Logger != "slog" {
		t.Errorf("Logger: got %q, want slog (global should survive)", got.Logger)
	}
}

func TestLocalFalseBoolOverridesGlobalTrue(t *testing.T) {
	dir := t.TempDir()
	global := writeYAML(t, dir, "global.yaml", `
docker: true
makefile: true
`)
	local := writeYAML(t, dir, "local.yaml", `
docker: false
`)

	got, err := loadFrom(global, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// local explicitly set docker to false — must win over global true
	if got.Docker == nil || *got.Docker {
		t.Errorf("Docker: expected false, got %v", got.Docker)
	}
	// makefile not touched by local — global true must survive
	if got.Makefile == nil || !*got.Makefile {
		t.Errorf("Makefile: expected true (from global), got %v", got.Makefile)
	}
}

func TestPartialLocalMerge(t *testing.T) {
	dir := t.TempDir()
	global := writeYAML(t, dir, "global.yaml", `
module_prefix: github.com/global-org
go_version: "1.24.0"
framework: gin
logger: slog
viper: true
docker: true
makefile: true
github: true
lint: true
`)
	local := writeYAML(t, dir, "local.yaml", `
logger: zap
swagger: true
`)

	got, err := loadFrom(global, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// local overrides
	if got.Logger != "zap" {
		t.Errorf("Logger: got %q, want zap", got.Logger)
	}
	if got.Swagger == nil || !*got.Swagger {
		t.Errorf("Swagger: expected true")
	}
	// global survives for untouched fields
	if got.ModulePrefix != "github.com/global-org" {
		t.Errorf("ModulePrefix: got %q, want github.com/global-org", got.ModulePrefix)
	}
	if got.Framework != "gin" {
		t.Errorf("Framework: got %q, want gin", got.Framework)
	}
	if got.GoVersion != "1.24.0" {
		t.Errorf("GoVersion: got %q, want 1.24.0", got.GoVersion)
	}
}

func TestInvalidYAMLReturnsError(t *testing.T) {
	dir := t.TempDir()
	bad := writeYAML(t, dir, "bad.yaml", `
framework: [unclosed bracket
`)

	_, err := loadFrom(bad, filepath.Join(dir, "local.yaml"))
	if err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestBoolVal(t *testing.T) {
	cases := []struct {
		name     string
		b        *bool
		fallback bool
		want     bool
	}{
		{"nil uses fallback true", nil, true, true},
		{"nil uses fallback false", nil, false, false},
		{"false pointer overrides true fallback", boolPtr(false), true, false},
		{"true pointer overrides false fallback", boolPtr(true), false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolVal(tc.b, tc.fallback)
			if got != tc.want {
				t.Errorf("BoolVal(%v, %v) = %v, want %v", tc.b, tc.fallback, got, tc.want)
			}
		})
	}
}

func TestMergeNilGlobal(t *testing.T) {
	local := &UserConfig{Framework: "echo", Docker: boolPtr(true)}
	got := merge(nil, local)
	if got.Framework != "echo" {
		t.Errorf("Framework: got %q, want echo", got.Framework)
	}
	if got.Docker == nil || !*got.Docker {
		t.Errorf("Docker: expected true")
	}
}

func TestMergeNilLocal(t *testing.T) {
	global := &UserConfig{Framework: "chi", Logger: "zerolog"}
	got := merge(global, nil)
	if got.Framework != "chi" {
		t.Errorf("Framework: got %q, want chi", got.Framework)
	}
	if got.Logger != "zerolog" {
		t.Errorf("Logger: got %q, want zerolog", got.Logger)
	}
}
