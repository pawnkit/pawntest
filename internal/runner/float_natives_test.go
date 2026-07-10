package runner

import (
	"math"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestFloatNativeHelpers(t *testing.T) {
	cell := func(value float32) backend.Cell { return backend.Cell(int32(math.Float32bits(value))) }

	value, err := callFloatNative("floatadd", nil, []backend.Cell{cell(1.5), cell(2.5)})
	if err != nil {
		t.Fatal(err)
	}

	if got := math.Float32frombits(uint32(value)); got != 4 {
		t.Fatalf("floatadd = %v, want 4", got)
	}

	value, err = callFloatNative("floatint", nil, []backend.Cell{cell(-2.75)})
	if err != nil || value != -2 {
		t.Fatalf("floatint = %d, %v; want -2", value, err)
	}
}

func TestTruncateMockOutputString(t *testing.T) {
	if got := truncateString("abcdef", 4); got != "abc" {
		t.Fatalf("truncateString() = %q, want abc", got)
	}

	if got := truncateString("hello", 0); got != "" {
		t.Fatalf("truncateString() = %q, want empty", got)
	}
}
