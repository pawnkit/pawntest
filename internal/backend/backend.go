package backend

import "github.com/pawnkit/goamx/vm"

type Cell = vm.Cell
type Public = vm.Public
type NativeFunc = vm.NativeFunc
type NativeContext = vm.NativeContext

type PublicCaller = vm.PublicCaller
type MemorySnapshotter = vm.MemorySnapshotter
type InstructionLimiter = vm.InstructionLimiter
type DebugLocator = vm.DebugLocator
type CoverageLocation = vm.CoverageLocation
type CoverageInstrumenter = vm.CoverageInstrumenter

type VM interface {
	Publics() ([]Public, error)
	Natives() ([]string, error)
	RegisterNative(name string, fn NativeFunc) error
	ExecPublic(index int, args ...Cell) (Cell, error)
	ExecMain(args ...Cell) (Cell, error)
	Reset() error
	Close() error
}

type Backend interface {
	LoadFile(path string) (VM, error)
	LoadBytes(name string, data []byte) (VM, error)
}

type GoAMXBackend struct{}

func NewGoAMXBackend() GoAMXBackend { return GoAMXBackend{} }

func (GoAMXBackend) LoadFile(path string) (VM, error) {
	return vm.LoadFile(path)
}

func (GoAMXBackend) LoadBytes(name string, data []byte) (VM, error) {
	return vm.LoadBytes(name, data)
}
