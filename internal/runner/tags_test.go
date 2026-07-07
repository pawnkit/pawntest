package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestTagExpressionsFilterPublics(t *testing.T) {
	t.Parallel()
	publics := []backend.Public{
		{Name: "test_fast"}, {Name: "test_slow"},
		{Name: "ptt_unit_fast"}, {Name: "ptt_unit_slow"}, {Name: "ptt_slow_slow"},
	}
	tests, err := (Runner{TagExpression: "unit & !slow"}).selectTests(publics)
	if err != nil {
		t.Fatal(err)
	}
	if len(tests) != 1 || tests[0].Name != "test_fast" {
		t.Fatalf("unexpected tests: %#v", tests)
	}
}

func TestInvalidTagExpression(t *testing.T) {
	t.Parallel()
	if _, err := parseTagExpression("unit & ("); err == nil {
		t.Fatal("expected invalid expression error")
	}
}
