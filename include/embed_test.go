package include

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"testing"
)

const pawnIdentifierLimit = 31

var guardPattern = regexp.MustCompile(`^#(?:if defined|define)\s+([A-Za-z_][A-Za-z0-9_]*)$`)
var nativePattern = regexp.MustCompile(`^native\s+(?:[A-Za-z_][A-Za-z0-9_]*:)?([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

func TestEmbeddedIncludesContainValidGuards(t *testing.T) {
	embedded := Files()
	if len(embedded) < 20 {
		t.Fatalf("embedded %d include files, want at least 20", len(embedded))
	}

	for _, path := range []string{"pawntest/scenario_assertions.inc", "pawntest/scenarios/http.inc", "pawntest/scenarios/database.inc"} {
		if len(embedded[path]) == 0 {
			t.Fatalf("nested include %s was not embedded", path)
		}
	}

	for path, data := range embedded {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "+#") {
				t.Fatalf("%s contains a malformed directive: %s", path, line)
			}

			match := guardPattern.FindStringSubmatch(line)
			if len(match) == 2 && len(match[1]) > pawnIdentifierLimit {
				t.Fatalf("%s guard %q has %d characters", path, match[1], len(match[1]))
			}

			match = nativePattern.FindStringSubmatch(line)
			if len(match) == 2 && len(match[1]) > pawnIdentifierLimit {
				t.Fatalf("%s native %q has %d characters", path, match[1], len(match[1]))
			}
		}

		if err := scanner.Err(); err != nil {
			t.Fatalf("scan %s: %v", path, err)
		}
	}
}
