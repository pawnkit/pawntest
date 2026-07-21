package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pawnkit/pawntest/internal/runner"
)

func TestPlainSummaryAndMessages(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{
		{Name: "test_pass", Status: runner.Pass},
		{Name: "test_fail", Status: runner.Fail, Message: "bad math", File: "math.pwn", Line: 12},
		{Name: "test_skip", Status: runner.Skip, Message: "later"},
		{Name: "test_error", Status: runner.Error, Message: "runtime"},
	}}

	var out bytes.Buffer
	if err := Plain(&out, suite); err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"PASS  test_pass",
		"FAIL  test_fail\n      bad math",
		"at math.pwn:12",
		"FAIL  4 tests across 1 file in 0ms",
		"1 passed, 1 failed, 1 skipped, 1 errored",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Plain output missing %q:\n%s", want, got)
		}
	}
}

func TestPlainColorizesStatusesAndSummary(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{
		{Name: "test_pass", Status: runner.Pass},
		{Name: "test_fail", Status: runner.Fail, Message: "bad math", File: "math.pwn", Line: 12},
		{Name: "test_skip", Status: runner.Skip, Message: "later"},
	}}

	var out bytes.Buffer
	if err := PlainWithOptions(&out, suite, PlainOptions{Color: true}); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"\x1b[32mPASS \x1b[0m test_pass",
		"\x1b[31mFAIL \x1b[0m test_fail",
		"\x1b[33mSKIP \x1b[0m test_skip",
		"\x1b[2mat math.pwn:12\x1b[0m",
		"\x1b[32m1 passed\x1b[0m, \x1b[31m1 failed\x1b[0m, \x1b[33m1 skipped\x1b[0m",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("colored Plain output missing %q:\n%q", want, out.String())
		}
	}
}

func TestPlainGroupsSourcesAndPrintsRerunCommand(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{
		{Name: "test_ok", Source: "tests/math.test.pwn", File: "tests/math.test.pwn", Status: runner.Pass},
		{Name: "test_addition", Source: "tests/math.test.pwn", File: "tests/math.test.pwn", Line: 12, Status: runner.Fail, Message: "expected: 5\nactual:   4"},
	}}

	var out bytes.Buffer
	if err := Plain(&out, suite); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"tests/math.test.pwn\n  PASS  test_ok",
		"      expected: 5\n        actual:   4",
		"at line 12",
		"rerun: pawntest 'tests/math.test.pwn' --run '^test_addition$'",
		"2 tests across 1 file",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("grouped output missing %q:\n%s", want, out.String())
		}
	}
}

func TestPlainQuietAndVerboseModes(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{
		{Name: "test_ok", Source: "math.test.pwn", Status: runner.Pass, Duration: 2 * time.Millisecond},
		{Name: "test_skip", Source: "math.test.pwn", Status: runner.Skip, Message: "later"},
		{Name: "test_fail", Source: "math.test.pwn", File: "math.test.pwn", Line: 4, Status: runner.Fail, Message: "bad", Duration: 3 * time.Millisecond},
	}}

	var quiet, verbose bytes.Buffer
	if err := PlainWithOptions(&quiet, suite, PlainOptions{Quiet: true}); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(quiet.String(), "test_ok") || strings.Contains(quiet.String(), "test_skip") || !strings.Contains(quiet.String(), "test_fail") {
		t.Fatalf("unexpected quiet output:\n%s", quiet.String())
	}

	if err := PlainWithOptions(&verbose, suite, PlainOptions{Verbose: true}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(verbose.String(), "test_ok  2ms") || !strings.Contains(verbose.String(), "test_fail  3ms") {
		t.Fatalf("verbose output missing durations:\n%s", verbose.String())
	}
}

func TestTAPSkipDirective(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{{Name: "test_skip", Status: runner.Skip, Message: "later"}}}

	var out bytes.Buffer
	if err := TAP(&out, suite); err != nil {
		t.Fatal(err)
	}

	if got, want := out.String(), "ok 1 - test_skip # SKIP later"; !strings.Contains(got, want) {
		t.Fatalf("TAP output missing %q:\n%s", want, got)
	}
}

func TestJSONUsesDurationMSAndSummary(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{{Name: "test_pass", Source: "tests/pass.test.pwn", Status: runner.Pass, Duration: 2 * time.Millisecond}}}

	var out bytes.Buffer
	if err := JSON(&out, suite); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{`"summary"`, `"total": 1`, `"duration_ms": 2`, `"source": "tests/pass.test.pwn"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("JSON output missing %q:\n%s", want, out.String())
		}
	}
}

func TestListJSONUsesDiscoverySchema(t *testing.T) {
	var out bytes.Buffer
	if err := ListJSON(&out, []DiscoveryTest{{Name: "test_addition", File: "tests/math.test.pwn", Line: 4}}); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{`"schemaVersion": 1`, `"id": "test_addition"`, `"label": "addition"`, `"file": "tests/math.test.pwn"`, `"line": 4`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("discovery output missing %q:\n%s", want, out.String())
		}
	}
}

func TestJUnitUsesSecondsDuration(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{{Name: "test_pass", Source: "tests/pass.test.pwn", Status: runner.Pass, Duration: 1500 * time.Millisecond}}}

	var out bytes.Buffer
	if err := JUnit(&out, suite); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), `time="1.500"`) {
		t.Fatalf("JUnit output missing seconds duration:\n%s", out.String())
	}

	for _, want := range []string{`classname="tests/pass.test.pwn"`, `file="tests/pass.test.pwn"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("JUnit output missing %q:\n%s", want, out.String())
		}
	}
}

func TestTAPIncludesSource(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{{Name: "test_pass", Source: "tests/pass.test.pwn", Status: runner.Pass}}}

	var out bytes.Buffer
	if err := TAP(&out, suite); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), `source: "tests/pass.test.pwn"`) {
		t.Fatalf("TAP output missing source:\n%s", out.String())
	}
}

func TestExpectedFailureStatusesAcrossReports(t *testing.T) {
	suite := runner.Suite{Results: []runner.Result{
		{Name: "test_known", Status: runner.XFail, Message: "known defect"},
		{Name: "test_stale", Status: runner.XPass, Message: "unexpected pass"},
	}}
	if !suite.Failed() {
		t.Fatal("xpass must fail the suite")
	}

	if summary := suite.Summary(); summary.XFailed != 1 || summary.XPassed != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}

	var plain, tap, junit bytes.Buffer
	if err := Plain(&plain, suite); err != nil {
		t.Fatal(err)
	}

	if err := TAP(&tap, suite); err != nil {
		t.Fatal(err)
	}

	if err := JUnit(&junit, suite); err != nil {
		t.Fatal(err)
	}

	for output, want := range map[string]string{
		plain.String(): "1 xfailed, 1 xpassed",
		tap.String():   "ok 1 - test_known # TODO expected failure",
		junit.String(): `failures="1"`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}
