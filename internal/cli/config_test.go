package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "pawntest.json")

	data := []byte(`{"pawncc":"./pawncc","include":["include"],"tests":["tests/..."],"format":"json","cache_dir":".cache","allow_unknown_natives":true,"allow_empty":true}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.PawnCC != "./pawncc" || cfg.Format != FormatJSON || !cfg.AllowUnknownNatives || !cfg.AllowEmpty {
		t.Fatalf("unexpected config: %#v", cfg)
	}

	if len(cfg.Tests) != 1 || cfg.Tests[0] != "tests/..." {
		t.Fatalf("tests = %#v, want tests/...", cfg.Tests)
	}
}

func TestLoadYAMLConfig(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "pawntest.yaml")

	data := []byte("pawncc: ./pawncc\ninclude:\n  - include\ntests:\n  - tests/...\nformat: tap\ncache_dir: .cache\nallow_unknown_natives: true\nallow_empty: true\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.PawnCC != "./pawncc" || cfg.Format != FormatTAP || !cfg.AllowUnknownNatives || !cfg.AllowEmpty {
		t.Fatalf("unexpected config: %#v", cfg)
	}

	if len(cfg.Include) != 1 || cfg.Include[0] != "include" {
		t.Fatalf("include = %#v, want include", cfg.Include)
	}
}

func TestLoadTOMLConfig(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "pawntest.toml")

	data := []byte("pawncc = \"./pawncc\"\ninclude = [\"include\"]\ntests = [\"tests/...\"]\nformat = \"tap\"\nallow_unknown_natives = true\nallow_empty = true\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.PawnCC != "./pawncc" || cfg.Format != FormatTAP || !cfg.AllowUnknownNatives || !cfg.AllowEmpty {
		t.Fatalf("unexpected config: %#v", cfg)
	}

	if len(cfg.Include) != 1 || cfg.Include[0] != "include" {
		t.Fatalf("include = %#v, want include", cfg.Include)
	}
}

func TestLoadConfigRejectsInvalidFormat(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "pawntest.json")
	if err := os.WriteFile(path, []byte(`{"format":"xmlish"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestApplyConfigDoesNotOverwriteExplicitValues(t *testing.T) {
	t.Parallel()

	cmd := TestCmd{
		Paths:       []string{"explicit.amx"},
		SharedFlags: SharedFlags{PawnCC: "pawncc-explicit", Include: []string{"explicit-include"}},
		Format:      FormatTAP,
	}
	cmd.applyConfig(Config{
		PawnCC:  "pawncc-config",
		Include: []string{"config-include"},
		Tests:   []string{"config-tests"},
		Format:  FormatJSON,
	})

	if cmd.Paths[0] != "explicit.amx" || cmd.PawnCC != "pawncc-explicit" || cmd.Include[0] != "explicit-include" || cmd.Format != FormatTAP {
		t.Fatalf("config overwrote explicit values: %#v", cmd)
	}
}
