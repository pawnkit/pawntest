package pawntest

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewRunnerUsesSafeDefaults(t *testing.T) {
	runner := NewRunner()
	if runner.MaxInstructions != DefaultMaxInstructions || runner.Isolation != "test" || runner.Repeat != 1 {
		t.Fatalf("NewRunner() = %#v", runner)
	}

	if got := (Runner{}).internal().MaxInstructions; got != DefaultMaxInstructions {
		t.Fatalf("zero-value instruction limit = %d", got)
	}
}

func TestRunnerContextCancellationStopsCompilation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := (Runner{}).RunFileContext(ctx, filepath.Join(t.TempDir(), "missing.test.pwn"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context cancellation", err)
	}
}

func TestRunnerCompilesSourceBeforeListing(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "sample.test.pwn")
	if err := os.WriteFile(src, []byte("#include <pawntest>\nTEST(sample) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	pawncc := fakePawnCC(t, dir)
	r := Runner{PawnCC: pawncc, CacheDir: filepath.Join(dir, "cache")}

	if _, err := r.List(src); err == nil {
		t.Fatal("expected fake AMX load error after source compilation")
	}

	if _, err := os.Stat(filepath.Join(dir, "cache", "include", "pawntest.inc")); err != nil {
		t.Fatal(err)
	}
}

func fakePawnCC(t *testing.T, dir string) string {
	t.Helper()

	path := filepath.Join(dir, "pawncc")
	if runtime.GOOS == "windows" {
		path += ".bat"
		script := "@echo off\r\n"
		script += "setlocal EnableDelayedExpansion\r\n"
		script += ":loop\r\n"
		script += "if \"%~1\"==\"\" goto done\r\n"
		script += "set \"arg=%~1\"\r\n"
		script += "if \"!arg:~0,2!\"==\"-o\" set \"out=!arg:~2!\"\r\n"
		script += "shift\r\n"
		script += "goto loop\r\n"
		script += ":done\r\n"
		script += "for %%I in (\"!out!\") do if not exist \"%%~dpI\" mkdir \"%%~dpI\"\r\n"

		script += "> \"!out!\" echo amx\r\n"
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}

		return path
	}

	script := "#!/bin/sh\n"
	script += "for arg in \"$@\"; do case \"$arg\" in -o*) out=${arg#-o};; esac; done\n"
	script += "mkdir -p \"$(dirname \"$out\")\"\n"

	script += "printf 'amx' > \"$out\"\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	return path
}
