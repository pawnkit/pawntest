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
	File  string `json:"file,omitempty"`
	Line  int    `json:"line,omitempty"`
}

type DiscoveryTest struct {
	Name string
	File string
	Line int
}

func ListJSON(w io.Writer, discovered []DiscoveryTest) error {
	tests := make([]discoveryTest, 0, len(discovered))
	for _, test := range discovered {
		tests = append(tests, discoveryTest{
			ID: test.Name, Label: strings.TrimPrefix(test.Name, "test_"), File: test.File, Line: test.Line,
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(discoveryDocument{SchemaVersion: 1, Tests: tests})
}
