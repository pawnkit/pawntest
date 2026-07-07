package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestSelectTestsSupportsRegexAndShuffle(t *testing.T) {
	t.Parallel()
	publics := []backend.Public{
		{Index: 1, Name: "test_alpha"},
		{Index: 2, Name: "test_beta"},
		{Index: 3, Name: "test_gamma"},
	}
	runner := Runner{Run: "alpha|gamma", Shuffle: true, Seed: 42, Repeat: 2}

	got, err := runner.selectTests(publics)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("selected %d tests, want 2", len(got))
	}
	runs := runner.testRuns(got)
	if len(runs) != 4 || runs[0].name != got[0].Name+" [attempt 1/2]" || runs[2].name != got[0].Name+" [attempt 2/2]" {
		t.Fatalf("unexpected repeated runs: %#v", runs)
	}
}

func TestSelectTestsRejectsInvalidRegex(t *testing.T) {
	t.Parallel()
	_, err := (Runner{Run: "["}).selectTests(nil)
	if err == nil {
		t.Fatal("expected invalid regex error")
	}
}
