//go:build integration

package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pawnkit/pawntest/internal/compiler"
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
		"--provider", filepath.Join(testsDir, "providers", "inventory.provider.pwn"),
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
		"PASS  test_player_event_flow",
		"PASS  test_vehicle_event_flow",
		"PASS  test_combat_event_flow",
		"PASS  test_movement_and_stream_events",
		"PASS  test_scenario_authoring",
		"PASS  test_strict_http_scenario",
		"PASS  test_provider_keeps_state",
		"PASS  test_provider_state_is_isolated",
		"PASS  test_provider_marshals_values",
		"PASS  test_provider_calls_test_public",
		"PASS  test_mock_overrides_provider",
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

	var multiTagOut, multiTagErr bytes.Buffer
	code = Run([]string{
		"--pawncc", pawncc,
		"--cache-dir", cacheDir,
		"--tags", "unit & !slow",
		filepath.Join(testsDir, "passing.test.pwn"),
		filepath.Join(testsDir, "cases.test.pwn"),
	}, &multiTagOut, &multiTagErr)
	if code != ExitOK || !strings.Contains(multiTagOut.String(), "PASS  test_even_two") {
		t.Fatalf("multi-file tag run exit=%d stderr=%q stdout=%q", code, multiTagErr.String(), multiTagOut.String())
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

func TestPawnCCIntegrationRunsIsolatedPlugin(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("x64 Linux fixture")
	}
	pawncc := os.Getenv("PAWNTEST_PAWNCC")
	if pawncc == "" && os.Getenv("PAWNTEST_PLUGIN_INTEGRATION") != "" {
		installed, err := compiler.InstallOpenMPCompiler(context.Background(), t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		pawncc = installed.Path
	}
	if pawncc == "" {
		t.Skip("set PAWNTEST_PAWNCC or PAWNTEST_PLUGIN_INTEGRATION to run plugin integration test")
	}
	compiler, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("C compiler unavailable")
	}
	pluginRoot, err := filepath.Abs(filepath.Join("..", "..", "..", "pawn-plugin-host"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	worker := filepath.Join(directory, "pawn-plugin-host-x64")
	loader := filepath.Join(directory, "pawn-plugin-loader")
	plugin := filepath.Join(directory, "fixture.so")
	commands := []*exec.Cmd{exec.Command("go", "build", "-o", worker, filepath.Join(pluginRoot, "cmd", "pawn-plugin-host-x64")), exec.Command(compiler, filepath.Join(pluginRoot, "native", "legacy_loader.c"), "-ldl", "-o", loader), exec.Command(compiler, "-shared", "-fPIC", filepath.Join(pluginRoot, "fixtures", "legacy_plugin.c"), "-o", plugin)}
	for _, command := range commands {
		command.Dir = pluginRoot
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("build fixture: %v: %s", err, output)
		}
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{"--pawncc", pawncc, "--cache-dir", filepath.Join(directory, "cache"), "--native-plugin", plugin, "--plugin-architecture", "x64", "--plugin-worker-64", worker, filepath.Join("..", "..", "testdata", "pawn", "plugin.test.pwn")}, &stdout, &stderr)
	if code != ExitOK || !strings.Contains(stdout.String(), "PASS  test_isolated_plugin_native") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}
