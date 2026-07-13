package runner

import (
	"slices"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestStrictScenarioFailuresReportResources(t *testing.T) {
	registry := newScenarioRegistry()
	registry.Database.connections[1] = &stubDatabaseConnection{}
	registry.Database.results[1] = &databaseResult{}
	registry.HTTP.responses[httpResponseKey{method: httpGet, url: "example.test"}] = []httpResponse{{code: 200}}
	registry.HTTP.unmatched = []httpRequest{{method: httpPost, url: "example.test"}}

	failures := registry.StrictFailures()

	want := []string{
		"unfreed database results: 1",
		"unclosed database connections: 1",
		`unconfigured HTTP request: method 2, URL "example.test"`,
		"unused HTTP responses: 1",
	}
	if !slices.Equal(failures, want) {
		t.Fatalf("strict failures = %v, want %v", failures, want)
	}
}

func TestStrictScenarioMarkerDetection(t *testing.T) {
	publics := []backend.Public{{Name: "test_example"}, {Name: "__pawntest_strict_scenarios"}}
	if !hasPublic(publics, "__pawntest_strict_scenarios") {
		t.Fatal("strict scenario marker was not detected")
	}
}
