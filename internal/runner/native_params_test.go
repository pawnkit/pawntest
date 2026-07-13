package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestNativeParamsValidatesAndReadsArguments(t *testing.T) {
	vm := &mockVM{strings: map[backend.Cell]string{100: "value"}}

	params := readNativeParams(vm, []backend.Cell{7, 100})
	if err := params.Require(2, "example"); err != nil {
		t.Fatal(err)
	}

	if params.Int(0) != 7 || params.Cell(1) != 100 {
		t.Fatal("native cells were not read")
	}

	value, err := params.String(1)
	if err != nil || value != "value" {
		t.Fatalf("string = %q, error = %v", value, err)
	}

	if err := params.Require(3, "example"); err == nil {
		t.Fatal("missing native argument was accepted")
	}
}
