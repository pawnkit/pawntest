package runner

import (
	"errors"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestRunnerFailFastStopsAfterError(t *testing.T) {
	vm := &failFastVM{}
	r := Runner{Backend: failFastBackend{vm: vm}, FailFast: true}

	suite, err := r.RunFile("failfast.amx")
	if err != nil {
		t.Fatal(err)
	}
	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}
	if suite.Results[0].Status != Error {
		t.Fatalf("status = %s, want %s", suite.Results[0].Status, Error)
	}
	if vm.executed != 1 {
		t.Fatalf("executed tests = %d, want 1", vm.executed)
	}
}

func TestRunnerReturnsTypedNoTestsError(t *testing.T) {
	vm := &failFastVM{publics: []backend.Public{{Index: 0, Name: markerPublic}}}
	r := Runner{Backend: failFastBackend{vm: vm}}

	_, err := r.RunFile("empty.amx")
	if !errors.Is(err, ErrNoTestsFound) {
		t.Fatalf("RunFile() error = %v, want ErrNoTestsFound", err)
	}
}

type failFastBackend struct {
	vm *failFastVM
}

func (b failFastBackend) LoadFile(path string) (backend.VM, error) {
	return b.vm, nil
}

func (b failFastBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	return b.vm, nil
}

type failFastVM struct {
	publics  []backend.Public
	executed int
}

func (vm *failFastVM) Publics() ([]backend.Public, error) {
	if vm.publics != nil {
		return vm.publics, nil
	}
	return []backend.Public{
		{Index: 0, Name: markerPublic},
		{Index: 1, Name: "test_errors"},
		{Index: 2, Name: "test_should_not_run"},
	}, nil
}

func (vm *failFastVM) Natives() ([]string, error) {
	return nil, nil
}

func (vm *failFastVM) RegisterNative(name string, fn backend.NativeFunc) error {
	return nil
}

func (vm *failFastVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	if index != 0 {
		vm.executed++
	}
	if index == 1 {
		return 0, errors.New("boom")
	}
	return 0, nil
}

func (vm *failFastVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *failFastVM) Reset() error {
	return nil
}

func (vm *failFastVM) Close() error {
	return nil
}

func (vm *failFastVM) ReadString(addr backend.Cell) (string, error) {
	return "", nil
}

func (vm *failFastVM) WriteString(addr backend.Cell, value string) error {
	return nil
}

func (vm *failFastVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *failFastVM) WriteCell(addr backend.Cell, value backend.Cell) error {
	return nil
}
