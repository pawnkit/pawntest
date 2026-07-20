package runner

import (
	"encoding/json"
	"io"
	"sort"
	"sync"

	"github.com/pawnkit/pawntest/internal/backend"
)

type Profile struct {
	mu           sync.Mutex
	instructions map[string]uint64
}

func NewProfile() *Profile { return &Profile{instructions: map[string]uint64{}} }

func (profile *Profile) observe(vm backend.VM, cip backend.Cell) {
	name := "<unknown>"

	if locator, ok := vm.(backend.DebugLocator); ok {
		_, _, function, found := locator.DebugLocation(cip)
		if found && function != "" {
			name = function
		}
	}

	profile.mu.Lock()
	profile.instructions[name]++
	profile.mu.Unlock()
}

func (profile *Profile) WriteJSON(writer io.Writer) error {
	profile.mu.Lock()
	defer profile.mu.Unlock()

	type row struct {
		Function     string `json:"function"`
		Instructions uint64 `json:"instructions"`
	}

	rows := make([]row, 0, len(profile.instructions))
	for name, count := range profile.instructions {
		rows = append(rows, row{name, count})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Instructions == rows[j].Instructions {
			return rows[i].Function < rows[j].Function
		}

		return rows[i].Instructions > rows[j].Instructions
	})

	return json.NewEncoder(writer).Encode(map[string]any{"version": 1, "functions": rows})
}
