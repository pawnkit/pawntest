package cli

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadProjectConfig(t *testing.T) {
	dir := t.TempDir()

	manifest := `{
  "pawnkit": {
    "schemaVersion": 1,
    "includePaths": ["include"],
    "tests": {"paths": ["tests/..."]},
    "tool": {"pawntest": {"format": "json", "jobs": 2}}
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "pawn.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(dir, "tests")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg, root, err := loadProjectConfig(nested)
	if err != nil {
		t.Fatal(err)
	}

	if root != dir || cfg.Format != FormatJSON || cfg.Jobs != 2 {
		t.Fatalf("root/config = %q %#v", root, cfg)
	}

	if want := filepath.Join(dir, "include"); len(cfg.Include) != 1 || cfg.Include[0] != want {
		t.Fatalf("include = %#v, want %q", cfg.Include, want)
	}

	if want := filepath.Join(dir, "tests", "..."); len(cfg.Tests) != 1 || cfg.Tests[0] != want {
		t.Fatalf("tests = %#v, want %q", cfg.Tests, want)
	}
}

func TestLoadProjectConfigOutsideProject(t *testing.T) {
	cfg, root, err := loadProjectConfig(t.TempDir())
	if err != nil || root != "" || !reflect.DeepEqual(cfg, Config{}) {
		t.Fatalf("got %#v %q %v", cfg, root, err)
	}
}

func TestLoadProjectConfigRejectsUnknownToolFields(t *testing.T) {
	dir := t.TempDir()
	manifest := `{"pawnkit":{"schemaVersion":1,"tool":{"pawntest":{"formt":"json"}}}}`

	if err := os.WriteFile(filepath.Join(dir, "pawn.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, _, err := loadProjectConfig(dir); err == nil {
		t.Fatal("unknown pawntest field was accepted")
	}
}
