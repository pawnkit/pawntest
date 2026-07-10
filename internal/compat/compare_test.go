package compat

import (
	"encoding/binary"
	"testing"

	goamx "github.com/pawnkit/goamx/vm"
	"github.com/pawnkit/pawntest/internal/backend"
)

func TestCompareEqualBackends(t *testing.T) {
	result, err := Compare("left", backend.NewGoAMXBackend(), "right", backend.NewGoAMXBackend(), Case{
		Name: "const.amx",
		AMX:  compatAMX(t, compatInstr(nil, goamx.OP_CONST_PRI, 42), compatInstr(nil, goamx.OP_HALT, 0)),
	})
	if err != nil {
		t.Fatal(err)
	}

	if !result.Equal() {
		t.Fatalf("expected equal comparison, got %#v", result)
	}
}

func compatAMX(t *testing.T, chunks ...[]byte) []byte {
	t.Helper()

	const headerSize = 56

	publics := uint32(headerSize)
	natives := publics + 8
	libs := natives
	data := make([]byte, libs)
	data = append(data, 31, 0)
	appendName := func(stubNameOffset uint32, name string) {
		off := uint32(len(data))
		binary.LittleEndian.PutUint32(data[stubNameOffset:stubNameOffset+4], off)
		data = append(data, name...)
		data = append(data, 0)
	}
	appendName(publics+4, "test_const")

	for len(data)%4 != 0 {
		data = append(data, 0)
	}

	cod := uint32(len(data))
	for _, chunk := range chunks {
		data = append(data, chunk...)
	}

	dat := uint32(len(data))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(data)))
	binary.LittleEndian.PutUint16(data[4:6], 0xf1e0)
	data[6] = 8
	data[7] = 11
	binary.LittleEndian.PutUint16(data[10:12], 8)
	binary.LittleEndian.PutUint32(data[12:16], cod)
	binary.LittleEndian.PutUint32(data[16:20], dat)
	binary.LittleEndian.PutUint32(data[20:24], dat)
	binary.LittleEndian.PutUint32(data[24:28], dat+256)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	binary.LittleEndian.PutUint32(data[40:44], libs)
	binary.LittleEndian.PutUint32(data[44:48], libs)
	binary.LittleEndian.PutUint32(data[48:52], libs)
	binary.LittleEndian.PutUint32(data[52:56], libs)
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(data)))

	return data
}

func compatInstr(data []byte, op goamx.Opcode, params ...backend.Cell) []byte {
	data = compatCell(data, backend.Cell(op))
	for _, param := range params {
		data = compatCell(data, param)
	}

	return data
}

func compatCell(data []byte, value backend.Cell) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(int32(value)))

	return append(data, buf[:]...)
}
