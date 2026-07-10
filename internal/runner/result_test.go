package runner

import (
	"strings"
	"testing"
)

func TestMergePhasePreservesMultipleFailures(t *testing.T) {
	t.Parallel()

	result := Result{Name: "test_value", Status: Pass}
	result = mergePhase(result, "test", Result{Status: Fail, Message: "body failed", File: "test.pwn", Line: 4})

	result = mergePhase(result, "teardown", Result{Status: Error, Message: "cleanup failed", File: "test.pwn", Line: 9})
	if result.Status != Error {
		t.Fatalf("status = %s, want error", result.Status)
	}

	if !strings.Contains(result.Message, "body failed") || !strings.Contains(result.Message, "teardown: cleanup failed") {
		t.Fatalf("messages were not aggregated: %q", result.Message)
	}

	if result.Line != 4 {
		t.Fatalf("primary failure line = %d, want 4", result.Line)
	}
}
