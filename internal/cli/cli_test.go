package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/pawnkit/pawntest/internal/compiler"
	"github.com/pawnkit/pawntest/internal/runner"
)

func TestDiscoveryPathsRecursive(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "math_test.amx")
	cmd := TestCmd{Paths: []string{dir, file}, Recursive: true}

	got := cmd.discoveryPaths()

	want := []string{dir + "/...", file}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("discoveryPaths mismatch (-want +got):\n%s", diff)
	}
}

func TestDoctorCommandReportsEnvironmentWithoutCompiler(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cacheDir := t.TempDir()

	var stdout, stderr bytes.Buffer

	code := Run([]string{"doctor", "--cache-dir", cacheDir}, &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; stderr=%q stdout=%q", code, ExitOK, stderr.String(), stdout.String())
	}

	for _, want := range []string{
		"pawntest doctor",
		"config: none",
		"cache: " + cacheDir,
		"include: ",
		"pawncc: not found",
		"sample: skipped",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("doctor output missing %q:\n%s", want, stdout.String())
		}
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestWriteReportColorAlwaysOnlyAppliesToPlain(t *testing.T) {
	t.Parallel()

	suite := runner.Suite{Results: []runner.Result{{Name: "test_pass", Status: runner.Pass}}}
	cmd := TestCmd{Format: FormatPlain, Color: ColorAlways}

	var plain bytes.Buffer
	if err := cmd.writeReport(&plain, suite); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(plain.String(), "\x1b[32mPASS \x1b[0m") {
		t.Fatalf("plain report was not colorized:\n%q", plain.String())
	}

	cmd.Format = FormatJSON

	var json bytes.Buffer
	if err := cmd.writeReport(&json, suite); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(json.String(), "\x1b[") {
		t.Fatalf("json report contains color escapes:\n%q", json.String())
	}
}

func TestEnsureCompilerAvailableInstallsWhenUserConfirms(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cmd := TestCmd{
		SharedFlags: SharedFlags{CacheDir: t.TempDir()},
		canPromptFn: func() bool { return true },
		confirmFn:   func() bool { return true },
		installFn: func(ctx context.Context, dir string) (*compiler.Compiler, error) {
			if dir == "" {
				t.Fatal("install dir is empty")
			}

			return compiler.Bare(filepath.Join(dir, "tools", "openmp-compiler", "pawncc")), nil
		},
	}

	var stderr bytes.Buffer
	if err := cmd.ensureCompilerAvailable(context.Background(), []string{"math.test.pwn"}, &stderr); err != nil {
		t.Fatal(err)
	}

	if cmd.PawnCC == "" {
		t.Fatal("PawnCC was not set after install")
	}

	if cmd.compilerCache == nil || cmd.compilerCache.Path != cmd.PawnCC {
		t.Fatalf("compilerCache = %+v, want Path %q", cmd.compilerCache, cmd.PawnCC)
	}

	if !bytes.Contains(stderr.Bytes(), []byte("Downloading openmultiplayer compiler")) {
		t.Fatalf("stderr = %q, want download message", stderr.String())
	}
}

func TestEnsureCompilerAvailableUsesCachedCompilerBeforePrompting(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cacheDir := t.TempDir()
	versionDir := filepath.Join(cacheDir, "tools", "openmp-compiler", "3.10.11", runtime.GOOS)
	markerDir := filepath.Join(versionDir, "pawncc-bin")

	realBinDir := filepath.Join(versionDir, "archive", "bin")
	for _, d := range []string{markerDir, realBinDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	name := "pawncc"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	cachedPawnCC := filepath.Join(realBinDir, name)
	if err := os.WriteFile(cachedPawnCC, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(markerDir, "pawncc.path"), []byte(cachedPawnCC+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := TestCmd{
		SharedFlags: SharedFlags{CacheDir: cacheDir},
		canPromptFn: func() bool { t.Fatal("prompted despite cached compiler"); return false },
		installFn: func(context.Context, string) (*compiler.Compiler, error) {
			t.Fatal("installed despite cached compiler")
			return nil, nil
		},
	}
	if err := cmd.ensureCompilerAvailable(context.Background(), []string{"math.test.pwn"}, io.Discard); err != nil {
		t.Fatal(err)
	}

	if cmd.PawnCC != cachedPawnCC {
		t.Fatalf("PawnCC = %q, want cached compiler %q", cmd.PawnCC, cachedPawnCC)
	}

	if cmd.compilerCache == nil || len(cmd.compilerCache.LibDirs) == 0 {
		t.Fatalf("compilerCache = %+v, want lib dirs derived from %q", cmd.compilerCache, cachedPawnCC)
	}
}

func TestEnsureCompilerAvailableDoesNotPromptForAMXOnlyRuns(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	cmd := TestCmd{}
	if err := cmd.ensureCompilerAvailable(context.Background(), []string{"already.amx"}, io.Discard); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureCompilerAvailableReturnsNotFoundWhenNonInteractive(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cmd := TestCmd{
		SharedFlags: SharedFlags{CacheDir: t.TempDir()},
		canPromptFn: func() bool { return false },
	}

	err := cmd.ensureCompilerAvailable(context.Background(), []string{"math.test.pwn"}, io.Discard)
	if !errors.Is(err, compiler.ErrPawnCCNotFound) {
		t.Fatalf("error = %v, want ErrPawnCCNotFound", err)
	}
}
