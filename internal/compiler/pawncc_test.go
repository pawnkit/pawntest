package compiler

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func FuzzParseIncludes(f *testing.F) {
	f.Add("#include <pawntest>\n#include \"local.inc\"\n")
	f.Add("#include")
	f.Fuzz(func(t *testing.T, source string) {
		_ = parseIncludes(source)
	})
}

func TestCompileUsesCacheAndPassesArgs(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "math.test.pwn")
	if err := os.WriteFile(src, []byte("public test_addition() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	log := filepath.Join(dir, "calls.log")
	pawncc := fakePawnCC(t, dir, log)

	var diagnostics strings.Builder

	out, err := Compile(src, Options{
		Compiler:    Bare(pawncc),
		Includes:    []string{"include"},
		Defines:     []string{"TESTING"},
		ExtraArgs:   []string{"-O0"},
		OutDir:      filepath.Join(dir, "cache"),
		Diagnostics: &diagnostics,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(diagnostics.String(), "compiler warning") {
		t.Fatalf("diagnostics = %q", diagnostics.String())
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

func TestCompileSerializesConcurrentCacheWrites(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "math.test.pwn")
	if err := os.WriteFile(src, []byte("public test_addition() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	log := filepath.Join(dir, "calls.log")
	pawncc := fakePawnCC(t, dir, log)
	opts := Options{Compiler: Bare(pawncc), OutDir: filepath.Join(dir, "cache")}

	var wait sync.WaitGroup

	errorsFound := make(chan error, 8)

	for range 8 {
		wait.Go(func() {
			_, err := Compile(src, opts)
			errorsFound <- err
		})
	}

	wait.Wait()
	close(errorsFound)

	for err := range errorsFound {
		if err != nil {
			t.Fatal(err)
		}
	}

	calls, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}

	if got := normalizeNewlines(string(calls)); got != "call\n" {
		t.Fatalf("compiler calls = %q, want one", got)
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

func TestCompileKeyIncludesCompilerContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "math.test.pwn")
	pawncc := filepath.Join(dir, "pawncc")

	if err := os.WriteFile(src, []byte("public test_addition() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(pawncc, []byte("compiler-one"), 0o755); err != nil {
		t.Fatal(err)
	}

	first, err := compileKey(src, pawncc, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(pawncc, []byte("compiler-two"), 0o755); err != nil {
		t.Fatal(err)
	}

	second, err := compileKey(src, pawncc, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("compile key did not change after compiler content changed")
	}
}

func TestCompileKeyIncludesNestedRelativeInclude(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "math.test.pwn")

	nestedDir := filepath.Join(dir, "modules")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatal(err)
	}

	helper := filepath.Join(nestedDir, "helper.inc")
	value := filepath.Join(nestedDir, "value.inc")

	if err := os.WriteFile(src, []byte("#include \"modules/helper.inc\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(helper, []byte("#include \"value.inc\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(value, []byte("#define VALUE 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	first, err := compileKey(src, "missing-pawncc", Options{})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(value, []byte("#define VALUE 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	second, err := compileKey(src, "missing-pawncc", Options{})
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("compile key did not change after nested include changed")
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
		script += "echo compiler warning 1>&2\r\n"
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
	script += "printf 'compiler warning\\n' >&2\n"
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
