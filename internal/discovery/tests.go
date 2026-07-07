package discovery

import (
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestPublics(publics []backend.Public) []backend.Public {
	var out []backend.Public
	for _, pub := range publics {
		if !strings.HasPrefix(pub.Name, "test_") {
			continue
		}
		if isFixturePublic(pub.Name) {
			continue
		}
		out = append(out, pub)
	}
	return out
}

func isFixturePublic(name string) bool {
	switch name {
	case "test_setup", "test_teardown", "test_suite_setup", "test_suite_teardown":
		return true
	default:
		return false
	}
}
