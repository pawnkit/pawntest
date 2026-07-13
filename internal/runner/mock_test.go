package runner

import (
	"strings"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestRunnerMocksUnknownNatives(t *testing.T) {
	vm := &mockVM{
		natives: map[string]backend.NativeFunc{},
		strings: mockStrings(),
	}
	r := Runner{Backend: mockBackend{vm: vm}}

	suite, err := r.RunFile("mock.amx")
	if err != nil {
		t.Fatal(err)
	}

	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}

	if suite.Results[0].Status != Pass {
		t.Fatalf("status = %s, want %s: %s", suite.Results[0].Status, Pass, suite.Results[0].Message)
	}

	if vm.nativeReturn != 7 {
		t.Fatalf("mocked native return = %d, want 7", vm.nativeReturn)
	}
}

func TestRunnerErrorsOnUnmockedUnknownNativeByDefault(t *testing.T) {
	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: mockStrings(), callWithoutMock: true}
	r := Runner{Backend: mockBackend{vm: vm}}

	suite, err := r.RunFile("mock.amx")
	if err != nil {
		t.Fatal(err)
	}

	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}

	if suite.Results[0].Status != Error {
		t.Fatalf("status = %s, want %s", suite.Results[0].Status, Error)
	}

	if !strings.Contains(suite.Results[0].Message, "native external_native is not mocked") {
		t.Fatalf("message = %q, want unmocked native error", suite.Results[0].Message)
	}
}

func TestRunnerAllowsUnmockedUnknownNativeWhenEnabled(t *testing.T) {
	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: mockStrings(), callWithoutMock: true}
	r := Runner{Backend: mockBackend{vm: vm}, AllowUnknownNatives: true}

	suite, err := r.RunFile("mock.amx")
	if err != nil {
		t.Fatal(err)
	}

	if len(suite.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(suite.Results))
	}

	if suite.Results[0].Status != Pass {
		t.Fatalf("status = %s, want %s: %s", suite.Results[0].Status, Pass, suite.Results[0].Message)
	}

	if vm.nativeReturn != 0 {
		t.Fatalf("unmocked native return = %d, want 0", vm.nativeReturn)
	}

	if len(suite.Results[0].Warnings) != 1 || !strings.Contains(suite.Results[0].Warnings[0], "external_native") {
		t.Fatalf("warnings = %#v", suite.Results[0].Warnings)
	}
}

type mockBackend struct {
	vm *mockVM
}

func mockStrings() map[backend.Cell]string {
	return map[backend.Cell]string{
		100: "external_native",
		400: "mock.pwn",
	}
}

func (b mockBackend) LoadFile(path string) (backend.VM, error) {
	return b.vm, nil
}

func (b mockBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	return b.vm, nil
}

type mockVM struct {
	natives         map[string]backend.NativeFunc
	strings         map[backend.Cell]string
	nativeReturn    backend.Cell
	callWithoutMock bool
}

func (vm *mockVM) Publics() ([]backend.Public, error) {
	return []backend.Public{{Index: 0, Name: markerPublic}, {Index: 1, Name: "test_mock"}}, nil
}

func (vm *mockVM) Natives() ([]string, error) {
	return []string{"external_native", "__pawntest_mock_return", "__pawntest_expect_calls", "__pawntest_expect_arg"}, nil
}

func (vm *mockVM) RegisterNative(name string, fn backend.NativeFunc) error {
	vm.natives[name] = fn
	return nil
}

func (vm *mockVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	if _, err := vm.natives["__pawntest_expect_calls"](vm, []backend.Cell{100, 1, 1, 400, 9}); err != nil {
		return 0, err
	}

	if _, err := vm.natives["__pawntest_expect_arg"](vm, []backend.Cell{100, 0, 0, 123, 400, 10}); err != nil {
		return 0, err
	}

	if !vm.callWithoutMock {
		if _, err := vm.natives["__pawntest_mock_return"](vm, []backend.Cell{100, 7}); err != nil {
			return 0, err
		}
	}

	ret, err := vm.natives["external_native"](vm, []backend.Cell{123, -1, 700})
	if err != nil {
		return 0, err
	}

	vm.nativeReturn = ret

	return 0, nil
}

func (vm *mockVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *mockVM) Reset() error {
	return nil
}

func (vm *mockVM) Close() error {
	return nil
}

func (vm *mockVM) ReadString(addr backend.Cell) (string, error) {
	return vm.strings[addr], nil
}

func (vm *mockVM) WriteString(addr backend.Cell, value string) error {
	vm.strings[addr] = value
	return nil
}

func (vm *mockVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, nil
}

func (vm *mockVM) WriteCell(addr, value backend.Cell) error {
	return nil
}
