package discovery

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestTestPublicsSkipsFixtures(t *testing.T) {
	publics := []backend.Public{
		{Index: 0, Name: "test_setup"},
		{Index: 1, Name: "test_addition"},
		{Index: 2, Name: "helper"},
		{Index: 3, Name: "test_suite_teardown"},
		{Index: 4, Name: "test_strings"},
	}

	got := TestPublics(publics)
	want := []backend.Public{{Index: 1, Name: "test_addition"}, {Index: 4, Name: "test_strings"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("TestPublics mismatch (-want +got):\n%s", diff)
	}
}
