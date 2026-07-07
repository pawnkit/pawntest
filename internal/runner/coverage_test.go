package runner

import (
	"bytes"
	"strings"
	"testing"
)

func TestCoverageWriters(t *testing.T) {
	coverage := NewCoverage()
	coverage.ensure("test.pwn", 2)
	coverage.ensure("test.pwn", 3)
	coverage.files["test.pwn"][2] = 4

	var lcov, json bytes.Buffer
	if err := coverage.WriteLCOV(&lcov); err != nil {
		t.Fatal(err)
	}
	if err := coverage.WriteJSON(&json); err != nil {
		t.Fatal(err)
	}
	for output, want := range map[string]string{
		lcov.String(): "DA:2,4\nDA:3,0\nLF:2\nLH:1",
		json.String(): `"count": 4`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("coverage output missing %q:\n%s", want, output)
		}
	}
}
