package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/pawnkit/pawntest/internal/backend"
)

type Coverage struct {
	mu    sync.Mutex
	files map[string]map[int]int
}

func NewCoverage() *Coverage {
	return &Coverage{files: map[string]map[int]int{}}
}

func (coverage *Coverage) prepare(vm backend.VM) {
	instrumenter, ok := vm.(backend.CoverageInstrumenter)

	_, hasLocator := vm.(backend.DebugLocator)
	if !ok || !hasLocator {
		return
	}

	coverage.mu.Lock()
	for _, location := range instrumenter.CoverageLocations() {
		if coverageRelevant(location.File, location.Line, location.Function) {
			coverage.ensure(location.File, location.Line)
		}
	}
	coverage.mu.Unlock()
}

func (coverage *Coverage) observe(vm backend.VM, cip backend.Cell) {
	locator, ok := vm.(backend.DebugLocator)
	if !ok {
		return
	}

	file, line, _, ok := locator.DebugLocation(cip)
	if !ok || file == "" || line <= 0 {
		return
	}

	coverage.mu.Lock()
	if lines := coverage.files[file]; lines != nil {
		if _, registered := lines[line]; registered {
			lines[line]++
		}
	}
	coverage.mu.Unlock()
}

func coverageRelevant(file string, line int, function string) bool {
	if filepath.Base(file) == "pawntest.inc" {
		return false
	}

	if function == markerPublic || strings.HasPrefix(function, tagPublicPrefix) {
		return false
	}

	if source, err := os.ReadFile(file); err == nil {
		lineCount := bytes.Count(source, []byte{'\n'})
		if len(source) > 0 && source[len(source)-1] != '\n' {
			lineCount++
		}

		return line <= lineCount
	}

	return true
}

func (coverage *Coverage) ensure(file string, line int) {
	if coverage.files[file] == nil {
		coverage.files[file] = map[int]int{}
	}

	if _, exists := coverage.files[file][line]; !exists {
		coverage.files[file][line] = 0
	}
}

func (coverage *Coverage) WriteLCOV(writer io.Writer) error {
	coverage.mu.Lock()
	defer coverage.mu.Unlock()

	for _, file := range sortedCoverageFiles(coverage.files) {
		if _, err := fmt.Fprintf(writer, "TN:pawntest\nSF:%s\n", file); err != nil {
			return err
		}

		lines := sortedCoverageLines(coverage.files[file])
		hit := 0

		for _, line := range lines {
			count := coverage.files[file][line]
			if count > 0 {
				hit++
			}

			if _, err := fmt.Fprintf(writer, "DA:%d,%d\n", line, count); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(writer, "LF:%d\nLH:%d\nend_of_record\n", len(lines), hit); err != nil {
			return err
		}
	}

	return nil
}

func (coverage *Coverage) WriteJSON(writer io.Writer) error {
	coverage.mu.Lock()
	defer coverage.mu.Unlock()

	type lineCoverage struct {
		Line  int `json:"line"`
		Count int `json:"count"`
	}

	out := map[string][]lineCoverage{}

	for _, file := range sortedCoverageFiles(coverage.files) {
		for _, line := range sortedCoverageLines(coverage.files[file]) {
			out[file] = append(out[file], lineCoverage{Line: line, Count: coverage.files[file][line]})
		}
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	return encoder.Encode(out)
}

func sortedCoverageFiles(files map[string]map[int]int) []string {
	out := make([]string, 0, len(files))
	for file := range files {
		out = append(out, file)
	}

	slices.Sort(out)

	return out
}

func sortedCoverageLines(lines map[int]int) []int {
	out := make([]int, 0, len(lines))
	for line := range lines {
		out = append(out, line)
	}

	slices.Sort(out)

	return out
}
