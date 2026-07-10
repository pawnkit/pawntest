package runner

import (
	"errors"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestRunnerRequiresPawntestIncludeMarker(t *testing.T) {
	t.Parallel()

	vm := &markerVM{publics: []backend.Public{{Index: 0, Name: "test_body"}}}
	r := Runner{Backend: markerBackend{vm: vm}}

	if _, err := r.List("plain.amx"); !errors.Is(err, ErrMissingPawntestInclude) {
		t.Fatalf("List() error = %v, want ErrMissingPawntestInclude", err)
	}

	if _, err := r.RunFile("plain.amx"); !errors.Is(err, ErrMissingPawntestInclude) {
		t.Fatalf("RunFile() error = %v, want ErrMissingPawntestInclude", err)
	}
}

type markerBackend struct {
	vm *markerVM
}

func (b markerBackend) LoadFile(path string) (backend.VM, error) {
	return b.vm, nil
}

func (b markerBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	return b.vm, nil
}

type markerVM struct {
	publics []backend.Public
}

func (vm *markerVM) Publics() ([]backend.Public, error) {
	return vm.publics, nil
}

func (vm *markerVM) Natives() ([]string, error) {
	return nil, nil
}

func (vm *markerVM) RegisterNative(name string, fn backend.NativeFunc) error {
	return nil
}

func (vm *markerVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *markerVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *markerVM) Reset() error {
	return nil
}

func (vm *markerVM) Close() error {
	return nil
}

func (vm *markerVM) ReadString(addr backend.Cell) (string, error) {
	return "", nil
}

func (vm *markerVM) WriteString(addr backend.Cell, value string) error {
	return nil
}

func (vm *markerVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *markerVM) WriteCell(addr, value backend.Cell) error {
	return nil
}
