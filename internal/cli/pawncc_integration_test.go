//go:build integration

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPawnCCIntegrationListsAndRunsCompiledSource(t *testing.T) {
	pawncc := os.Getenv("PAWNTEST_PAWNCC")
	if pawncc == "" {
		t.Skip("set PAWNTEST_PAWNCC to run pawncc integration test")
	}
	testsDir := filepath.Join("..", "..", "testdata", "pawn")
	cacheDir := filepath.Join(t.TempDir(), "cache")

	var listOut, listErr bytes.Buffer
	code := Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		"--list",
		filepath.Join(testsDir, "passing.test.pwn"),
	}, &listOut, &listErr)
	if code != ExitOK {
		t.Fatalf("list exit = %d, want %d; stderr=%q stdout=%q", code, ExitOK, listErr.String(), listOut.String())
	}
	if !strings.Contains(listOut.String(), "test_addition") {
		t.Fatalf("list output = %q, want test_addition", listOut.String())
	}

	var runOut, runErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		filepath.Join(testsDir, "passing.test.pwn"),
	}, &runOut, &runErr)
	if code != ExitOK {
		t.Fatalf("run exit = %d, want %d; stderr=%q stdout=%q", code, ExitOK, runErr.String(), runOut.String())
	}
	if !strings.Contains(runOut.String(), "PASS  test_addition") {
		t.Fatalf("run output = %q, want pass test_addition", runOut.String())
	}

	var mockOut, mockErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		filepath.Join(testsDir, "natives.test.pwn"),
	}, &mockOut, &mockErr)
	if code != ExitOK {
		t.Fatalf("mock run exit = %d, want %d; stderr=%q stdout=%q", code, ExitOK, mockErr.String(), mockOut.String())
	}
	if !strings.Contains(mockOut.String(), "PASS  test_unknown_native_mock") {
		t.Fatalf("mock run output = %q, want pass test_unknown_native_mock", mockOut.String())
	}

	var suiteOut, suiteErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", filepath.Join(t.TempDir(), "suite-cache"),
		"--allow-unknown-natives",
		testsDir,
	}, &suiteOut, &suiteErr)
	if code != ExitFailed {
		t.Fatalf("suite exit = %d, want %d; stderr=%q stdout=%q", code, ExitFailed, suiteErr.String(), suiteOut.String())
	}
	for _, want := range []string{
		"PASS  test_control_flow_and_arrays",
		"XFAIL test_expected_failure",
		"XPASS test_unexpected_pass",
		"PASS  test_even_two",
		"PASS  test_even_four",
		"PASS  test_sum_case",
		"PASS  test_declarative_syntax",
		"PASS  compile:compiler_errors.test",
		"FAIL  test_failure",
		"PASS  test_fixtures_run_in_order",
		"PASS  test_named_fixture_is_available",
		"PASS  test_float_helpers",
		"PASS  test_rich_assertions",
		"PASS  test_mock_expectations",
		"PASS  test_include_dependency",
		"PASS  test_isolation_first",
		"PASS  test_isolation_second",
		"PASS  test_unknown_native_mock",
		"PASS  test_mock_return_sequence",
		"PASS  test_mock_output_parameters",
		"PASS  test_mock_callback",
		"PASS  test_virtual_time",
		"PASS  test_pending_callback_order",
		"PASS  test_string_snapshot",
		"PASS  test_integer_property",
		"PASS  test_player_scenario",
		"FAIL  test_property_shrinks_failure",
		"ERROR test_divide_by_zero_errors",
		"SKIP  test_skip_example",
	} {
		if !strings.Contains(suiteOut.String(), want) {
			t.Fatalf("suite output missing %q:\nstdout=%s\nstderr=%s", want, suiteOut.String(), suiteErr.String())
		}
	}

	var tagOut, tagErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		"--tags", "unit & !slow",
		filepath.Join(testsDir, "cases.test.pwn"),
	}, &tagOut, &tagErr)
	if code != ExitOK || !strings.Contains(tagOut.String(), "PASS  test_even_two") || strings.Contains(tagOut.String(), "test_expected_failure") {
		t.Fatalf("tag run exit=%d stderr=%q stdout=%q", code, tagErr.String(), tagOut.String())
	}

	coveragePath := filepath.Join(t.TempDir(), "coverage.lcov")
	var coverageOut, coverageErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		"--coverage",
		"--coverage-output", coveragePath,
		filepath.Join(testsDir, "passing.test.pwn"),
	}, &coverageOut, &coverageErr)
	if code != ExitOK {
		t.Fatalf("coverage run exit=%d stderr=%q stdout=%q", code, coverageErr.String(), coverageOut.String())
	}
	coverageData, err := os.ReadFile(coveragePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(coverageData, []byte("SF:"+filepath.Join(testsDir, "passing.test.pwn"))) || bytes.Contains(coverageData, []byte("pawntest.inc")) {
		t.Fatalf("unexpected coverage:\n%s", coverageData)
	}
}
