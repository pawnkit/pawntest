package report

import (
	"encoding/json"
	"io"
	"strings"
)

type discoveryDocument struct {
	SchemaVersion int             `json:"schemaVersion"`
	Tests         []discoveryTest `json:"tests"`
}

type discoveryTest struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func ListJSON(w io.Writer, names []string) error {
	tests := make([]discoveryTest, 0, len(names))
	for _, name := range names {
		tests = append(tests, discoveryTest{ID: name, Label: strings.TrimPrefix(name, "test_")})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(discoveryDocument{SchemaVersion: 1, Tests: tests})
}
