//go:build cgo && amx_c

package compat

/*
#cgo CFLAGS: -std=c99 -DPAWN_CELL_SIZE=32 -I${SRCDIR}/camx
#include <stdlib.h>
#include "c_bridge.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/pawnkit/pawntest/internal/backend"
)

type CBackend struct{}

func NewCBackend() *CBackend { return &CBackend{} }

func (b *CBackend) LoadFile(path string) (backend.VM, error) {
	return nil, errors.New("C compatibility backend only supports LoadBytes")
}

func (b *CBackend) LoadBytes(name string, data []byte) (backend.VM, error) {
	if len(data) == 0 {
		return nil, errors.New("empty AMX image")
	}
	var code C.int
	vm := C.pawntest_amx_load(unsafe.Pointer(&data[0]), C.size_t(len(data)), &code)
	if vm == nil {
		return nil, cAMXError(code)
	}
	return &cVM{vm: vm}, nil
}

type cVM struct{ vm *C.pawntest_amx }

func (vm *cVM) Publics() ([]backend.Public, error) {
	count := int(C.pawntest_amx_num_publics(vm.vm))
	if count < 0 {
		return nil, errors.New("could not read C AMX publics")
	}
	out := make([]backend.Public, 0, count)
	for i := 0; i < count; i++ {
		name := make([]byte, 256)
		if code := C.pawntest_amx_public_name(vm.vm, C.int(i), (*C.char)(unsafe.Pointer(&name[0])), C.size_t(len(name))); code != 0 {
			return nil, cAMXError(code)
		}
		out = append(out, backend.Public{Index: i, Name: C.GoString((*C.char)(unsafe.Pointer(&name[0])))})
	}
	return out, nil
}

func (vm *cVM) Natives() ([]string, error) {
	count := int(C.pawntest_amx_num_natives(vm.vm))
	if count < 0 {
		return nil, errors.New("could not read C AMX natives")
	}
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		name := make([]byte, 256)
		if code := C.pawntest_amx_native_name(vm.vm, C.int(i), (*C.char)(unsafe.Pointer(&name[0])), C.size_t(len(name))); code != 0 {
			return nil, cAMXError(code)
		}
		out = append(out, C.GoString((*C.char)(unsafe.Pointer(&name[0]))))
	}
	return out, nil
}

func (vm *cVM) RegisterNative(name string, fn backend.NativeFunc) error {
	return fmt.Errorf("C compatibility backend cannot register Go native %s", name)
}

func (vm *cVM) ExecPublic(index int, args ...backend.Cell) (backend.Cell, error) {
	var ptr *C.int32_t
	if len(args) > 0 {
		ptr = (*C.int32_t)(unsafe.Pointer(&args[0]))
	}
	var result C.int32_t
	code := C.pawntest_amx_exec(vm.vm, C.int(index), ptr, C.int(len(args)), &result)
	if code != 0 {
		return 0, cAMXError(code)
	}
	return backend.Cell(result), nil
}

func (vm *cVM) ExecMain(args ...backend.Cell) (backend.Cell, error) {
	return vm.ExecPublic(-1, args...)
}
func (vm *cVM) Reset() error { return nil }
func (vm *cVM) Close() error { C.pawntest_amx_free(vm.vm); vm.vm = nil; return nil }
func (vm *cVM) ReadString(addr backend.Cell) (string, error) {
	return "", errors.New("not exposed by C compatibility backend")
}

func (vm *cVM) WriteString(addr backend.Cell, value string) error {
	return errors.New("not exposed by C compatibility backend")
}

func (vm *cVM) ReadCell(addr backend.Cell) (backend.Cell, error) {
	return 0, errors.New("not exposed by C compatibility backend")
}

func (vm *cVM) WriteCell(addr, value backend.Cell) error {
	return errors.New("not exposed by C compatibility backend")
}

func cAMXError(code C.int) error {
	return fmt.Errorf("canonical C AMX: %s (code %d)", C.GoString(C.pawntest_amx_error(code)), int(code))
}
