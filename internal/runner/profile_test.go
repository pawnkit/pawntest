package runner

import (
	"bytes"
	"strings"
	"testing"
)

func TestProfileJSON(t *testing.T) {
	profile := NewProfile()
	profile.instructions["main"] = 7

	var output bytes.Buffer
	if err := profile.WriteJSON(&output); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), `"version":1`) || !strings.Contains(output.String(), `"instructions":7`) {
		t.Fatalf("output = %s", output.String())
	}
}
