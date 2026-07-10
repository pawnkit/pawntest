package compiler

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCompileUsesCacheAndPassesArgs(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "math.test.pwn")
	if err := os.WriteFile(src, []byte("public test_addition() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	log := filepath.Join(dir, "calls.log")
	pawncc := fakePawnCC(t, dir, log)

	out, err := Compile(src, Options{
		Compiler:  Bare(pawncc),
		Includes:  []string{"include"},
		Defines:   []string{"TESTING"},
		ExtraArgs: []string{"-O0"},
		OutDir:    filepath.Join(dir, "cache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if filepath.Base(out) == "math_test.amx" {
		t.Fatalf("Compile() output = %q, want cache-keyed filename", out)
	}

	second, err := Compile(src, Options{
		Compiler:  Bare(pawncc),
		Includes:  []string{"include"},
		Defines:   []string{"TESTING"},
		ExtraArgs: []string{"-O0"},
		OutDir:    filepath.Join(dir, "cache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if second != out {
		t.Fatalf("second Compile() = %q, want %q", second, out)
	}

	calls, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}

	if got := normalizeNewlines(string(calls)); got != "call\n" {
		t.Fatalf("fake pawncc calls = %q, want one call", got)
	}

	forced, err := Compile(src, Options{
		Compiler: Bare(pawncc),
		OutDir:   filepath.Join(dir, "cache"),
		Count:    1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if forced == "" {
		t.Fatal("forced Compile() returned empty output")
	}

	calls, err = os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}

	if got := normalizeNewlines(string(calls)); got != "call\ncall\n" {
		t.Fatalf("fake pawncc calls after forced compile = %q, want two calls", got)
	}
}

func TestBuildArgsDisablesCompactOutputAndEnablesDebugInfo(t *testing.T) {
	got := buildArgs("out.amx", "math.test.pwn", Options{
		Includes:  []string{"include"},
		Defines:   []string{"TESTING"},
		ExtraArgs: []string{"-O0"},
	})

	want := []string{"-C-", "-d2", "-oout.amx", "-iinclude", "-DTESTING", "-O0", "math.test.pwn"}
	if len(got) != len(want) {
		t.Fatalf("len(buildArgs()) = %d, want %d: %#v", len(got), len(want), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("buildArgs()[%d] = %q, want %q (all args: %#v)", i, got[i], want[i], got)
		}
	}
}

func TestCompileKeyIncludesUserIncludeContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "math.test.pwn")

	incDir := filepath.Join(dir, "include")
	if err := os.Mkdir(incDir, 0o755); err != nil {
		t.Fatal(err)
	}

	inc := filepath.Join(incDir, "helper.inc")

	if err := os.WriteFile(src, []byte("#include <helper>\npublic test_addition() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(inc, []byte("#define VALUE 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	first, err := compileKey(src, "pawncc", Options{Includes: []string{incDir}})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(inc, []byte("#define VALUE 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	second, err := compileKey(src, "pawncc", Options{Includes: []string{incDir}})
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("compile key did not change after included file content changed")
	}
}

func TestParseIncludes(t *testing.T) {
	got := parseIncludes("#include <pawntest>\n#include \"local.inc\"\n#include <nested/path>\n")

	want := []string{"pawntest", "local.inc", "nested/path"}
	if len(got) != len(want) {
		t.Fatalf("parseIncludes() = %#v, want %#v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parseIncludes()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func fakePawnCC(t *testing.T, dir, log string) string {
	t.Helper()

	path := filepath.Join(dir, "pawncc")
	if runtime.GOOS == "windows" {
		path += ".bat"
		script := "@echo off\r\n"
		script += "setlocal EnableDelayedExpansion\r\n"
		script += "echo call>> \"%PAWNTEST_CALL_LOG%\"\r\n"
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

		t.Setenv("PAWNTEST_CALL_LOG", log)

		return path
	}

	script := "#!/bin/sh\n"
	script += "printf 'call\\n' >> \"$PAWNTEST_CALL_LOG\"\n"
	script += "for arg in \"$@\"; do case \"$arg\" in -o*) out=${arg#-o};; esac; done\n"
	script += "mkdir -p \"$(dirname \"$out\")\"\n"

	script += "printf 'amx' > \"$out\"\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PAWNTEST_CALL_LOG", log)

	return path
}

func normalizeNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
