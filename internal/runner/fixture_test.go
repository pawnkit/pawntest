package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestRunnerExecutesSetupFixtureBeforeTest(t *testing.T) {
	vm := &fixtureVM{
		natives: map[string]backend.NativeFunc{},
		strings: map[backend.Cell]string{
			100: "setup failed",
			200: "fixture.pwn",
		},
	}
	r := Runner{Backend: fixtureBackend{vm: vm}}

	suite, err := r.RunFile("fixture.amx")
	if err != nil {
		t.Fatal(err)
	}

	if vm.bodyRan {
		t.Fatal("test body ran after setup failure")
	}

	if !vm.teardownRan {
		t.Fatal("teardown did not run after setup failure")
	}

	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}

	result := suite.Results[0]
	if result.Status != Fail {
		t.Fatalf("status = %s, want %s", result.Status, Fail)
	}

	if result.Message != "setup: setup failed" {
		t.Fatalf("message = %q, want setup phase failure", result.Message)
	}
}

type fixtureBackend struct {
	vm *fixtureVM
}

func (b fixtureBackend) LoadFile(path string) (backend.VM, error) {
	return b.vm, nil
}

func (b fixtureBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	return b.vm, nil
}

type fixtureVM struct {
	natives     map[string]backend.NativeFunc
	strings     map[backend.Cell]string
	bodyRan     bool
	teardownRan bool
}

func (vm *fixtureVM) Publics() ([]backend.Public, error) {
	return []backend.Public{
		{Index: 0, Name: markerPublic},
		{Index: 1, Name: "test_setup"},
		{Index: 2, Name: "test_body"},
		{Index: 3, Name: "test_teardown"},
	}, nil
}

func (vm *fixtureVM) Natives() ([]string, error) {
	return nil, nil
}

func (vm *fixtureVM) RegisterNative(name string, fn backend.NativeFunc) error {
	vm.natives[name] = fn
	return nil
}

func (vm *fixtureVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	if index == 1 {
		return vm.natives["__pawntest_fail"](vm, []backend.Cell{100, 200, 7})
	}

	if index == 3 {
		vm.teardownRan = true
		return 0, nil
	}

	vm.bodyRan = true

	return 0, nil
}

func (vm *fixtureVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *fixtureVM) Reset() error {
	return nil
}

func (vm *fixtureVM) Close() error {
	return nil
}

func (vm *fixtureVM) ReadString(addr backend.Cell) (string, error) {
	return vm.strings[addr], nil
}

func (vm *fixtureVM) WriteString(addr backend.Cell, value string) error {
	vm.strings[addr] = value
	return nil
}

func (vm *fixtureVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *fixtureVM) WriteCell(addr, value backend.Cell) error {
	return nil
}
