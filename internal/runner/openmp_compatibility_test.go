package runner

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestScenarioNativesExistInOpenMPSnapshot(t *testing.T) {
	data, err := os.ReadFile("testdata/openmp_natives.txt")
	if err != nil {
		t.Fatal(err)
	}

	official := map[string]bool{}

	for line := range strings.Lines(string(data)) {
		name := strings.TrimSpace(line)
		if name != "" && !strings.HasPrefix(name, "#") {
			official[name] = true
		}
	}

	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}
	registry := newScenarioRegistry()

	context := &executionContext{state: &nativeState{status: Pass}, mocks: newMockState(), scenarios: registry}
	if err := registry.Register(vm, context); err != nil {
		t.Fatal(err)
	}

	missing := []string{}

	for name := range vm.natives {
		if !isPawntestNative(name) && !official[name] {
			missing = append(missing, name)
		}
	}

	slices.Sort(missing)

	if len(missing) != 0 {
		t.Fatalf("scenario natives missing from open.mp snapshot: %s", strings.Join(missing, ", "))
	}
}
