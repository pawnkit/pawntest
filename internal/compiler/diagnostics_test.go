package compiler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiagnosticExpectations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.test.pwn")

	source := "#include <pawntest>\n// pawntest: expect-error 017\n// pawntest: expect-warning 213\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := DiagnosticExpectations(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got[0].Kind != "error" || got[0].Code != "017" || got[1].Kind != "warning" || got[1].Code != "213" {
		t.Fatalf("unexpected expectations: %#v", got)
	}
}

func TestDiagnosticExpectationsRequiresPawntestInclude(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.test.pwn")
	if err := os.WriteFile(path, []byte("// pawntest: expect-error 017\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := DiagnosticExpectations(path); err == nil {
		t.Fatal("expected missing pawntest include error")
	}
}
