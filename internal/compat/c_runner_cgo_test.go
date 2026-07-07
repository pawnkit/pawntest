//go:build cgo && amx_c

package compat

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
	goamx "github.com/pawnkit/goamx/vm"
)

func TestCanonicalCParity(t *testing.T) {
	cases := []struct {
		name string
		code []byte
	}{
		{"constant", compatCode(
			compatInstr(nil, goamx.OP_CONST_PRI, 42),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
		{"arithmetic", compatCode(
			compatInstr(nil, goamx.OP_CONST_PRI, 6),
			compatInstr(nil, goamx.OP_CONST_ALT, 7),
			compatInstr(nil, goamx.OP_SMUL),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
		{"comparison", compatCode(
			compatInstr(nil, goamx.OP_CONST_PRI, 7),
			compatInstr(nil, goamx.OP_CONST_ALT, 7),
			compatInstr(nil, goamx.OP_EQ),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
		{"memory", compatCode(
			compatInstr(nil, goamx.OP_CONST, 0, 41),
			compatInstr(nil, goamx.OP_INC, 0),
			compatInstr(nil, goamx.OP_LOAD_PRI, 0),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
		{"bounds-error", compatCode(
			compatInstr(nil, goamx.OP_CONST_PRI, 2),
			compatInstr(nil, goamx.OP_BOUNDS, 1),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
		{"divide-error", compatCode(
			compatInstr(nil, goamx.OP_CONST_PRI, 0),
			compatInstr(nil, goamx.OP_CONST_ALT, 1),
			compatInstr(nil, goamx.OP_SDIV),
			compatInstr(nil, goamx.OP_HALT, 0),
		)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Compare("go", backend.NewGoAMXBackend(), "canonical-c", NewCBackend(), Case{
				Name: tc.name + ".amx",
				AMX:  compatAMX(t, tc.code),
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.Equal() {
				t.Fatalf("parity mismatch: %#v", result)
			}
		})
	}
}

func compatCode(instructions ...[]byte) []byte {
	var code []byte
	for _, instruction := range instructions {
		code = append(code, instruction...)
	}
	return code
}

func TestCanonicalCParityCompiledCompactCorpus(t *testing.T) {
	pawncc := os.Getenv("PAWNTEST_PAWNCC")
	if pawncc == "" {
		t.Skip("set PAWNTEST_PAWNCC to run compiled C parity")
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "parity.pwn")
	program := `
forward test_add(a, b);
forward test_control(value);

public test_add(a, b)
{
    return a + b;
}

public test_control(value)
{
    new values[3] = {2, 4, 6};
    if (value > 0)
        return values[2] * 7;
    return -1;
}
`
	if err := os.WriteFile(source, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "parity.amx")
	cmd := exec.Command(pawncc, source, "-o"+out)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pawncc: %v\n%s", err, output)
	}
	image, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Compare("go", backend.NewGoAMXBackend(), "canonical-c", NewCBackend(), Case{
		Name: "compiled-compact.amx", AMX: image,
		PublicArgs: map[string][]backend.Cell{"test_add": {20, 22}, "test_control": {1}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equal() {
		t.Fatalf("compiled parity mismatch: %#v", result)
	}
}
