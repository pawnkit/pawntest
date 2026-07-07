package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestRunnerRecordsAssertionNativeFailure(t *testing.T) {
	vm := newFakeVM()
	r := Runner{Backend: fakeBackend{vm: vm}}

	suite, err := r.RunFile("test.amx")
	if err != nil {
		t.Fatal(err)
	}
	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}
	got := suite.Results[0]
	if got.Status != Fail {
		t.Fatalf("status = %s, want %s", got.Status, Fail)
	}
	if got.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestAssertNativePassesOnNonZeroCondition(t *testing.T) {
	state := &nativeState{}
	vm := newFakeVM()
	fn := nativeAssert(state)

	got, err := fn(vm, []backend.Cell{1, 100, 300, 12})
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("nativeAssert() = %d, want 1", got)
	}
	if state.status != "" || state.message != "" {
		t.Fatalf("nativeAssert() recorded failure state: %+v", state)
	}
}

type fakeBackend struct {
	vm *fakeVM
}

func (b fakeBackend) LoadFile(path string) (backend.VM, error) {
	return b.vm, nil
}

func (b fakeBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	return b.vm, nil
}

type fakeVM struct {
	natives map[string]backend.NativeFunc
	strings map[backend.Cell]string
}

func newFakeVM() *fakeVM {
	return &fakeVM{
		natives: map[string]backend.NativeFunc{},
		strings: map[backend.Cell]string{
			100: "actual",
			200: "expected",
			300: "test.pwn",
		},
	}
}

func (vm *fakeVM) Publics() ([]backend.Public, error) {
	return []backend.Public{{Index: 0, Name: markerPublic}, {Index: 1, Name: "test_addition"}}, nil
}

func (vm *fakeVM) Natives() ([]string, error) {
	return nil, nil
}

func (vm *fakeVM) RegisterNative(name string, fn backend.NativeFunc) error {
	vm.natives[name] = fn
	return nil
}

func (vm *fakeVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	return vm.natives["__pawntest_assert_eq"](vm, []backend.Cell{1, 2, 100, 200, 300, 12})
}

func (vm *fakeVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *fakeVM) Reset() error {
	return nil
}

func (vm *fakeVM) Close() error {
	return nil
}

func (vm *fakeVM) ReadString(addr backend.Cell) (string, error) {
	return vm.strings[addr], nil
}

func (vm *fakeVM) WriteString(addr backend.Cell, value string) error {
	vm.strings[addr] = value
	return nil
}

func (vm *fakeVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *fakeVM) WriteCell(addr backend.Cell, value backend.Cell) error {
	return nil
}
