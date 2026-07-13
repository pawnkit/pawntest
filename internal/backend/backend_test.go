package backend

import "testing"

func FuzzGoAMXLoadBytes(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("AMX"))
	f.Fuzz(func(t *testing.T, data []byte) {
		vm, err := (GoAMXBackend{}).LoadBytes("fuzz.amx", data)
		if err == nil {
			_ = vm.Close()
		}
	})
}
