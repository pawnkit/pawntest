package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateGeneratedSymbolsRejectsLongTestName(t *testing.T) {
	path := writeSymbolTest(t, `
// TEST(this_comment_does_not_count)
TEST(police_vehicle_models_are_recognized)
{
}
`)

	err := validateGeneratedSymbols(path)
	if err == nil {
		t.Fatal("validateGeneratedSymbols() returned nil")
	}

	for _, want := range []string{
		"math.test.pwn:3",
		`TEST name "police_vehicle_models_are_recognized"`,
		`"test_police_vehicle_models_are_recognized"`,
		"maximum source name length is 26",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err, want)
		}
	}
}

func TestValidateGeneratedSymbolsRejectsLongFixtureName(t *testing.T) {
	path := writeSymbolTest(t, "FIXTURE(a_fixture_name_that_is_too_long)\n")

	err := validateGeneratedSymbols(path)
	if err == nil || !strings.Contains(err.Error(), "maximum source name length is 14") {
		t.Fatalf("error = %v, want fixture length guidance", err)
	}
}

func TestValidateGeneratedSymbolsAcceptsSupportedMacros(t *testing.T) {
	path := writeSymbolTest(t, `
TEST(short_name) {}
TEST_CASE(short_case, callback, 1)
TEST_CASE2(short_case_two, callback, 1, 2)
TEST_CASE3(short_case_three, callback, 1, 2, 3)
FIXTURE(short_fixture) {}
`)

	if err := validateGeneratedSymbols(path); err != nil {
		t.Fatal(err)
	}
}

func writeSymbolTest(t *testing.T, source string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "math.test.pwn")
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}
