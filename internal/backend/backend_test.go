package backend

import "testing"

const maxFuzzImageSize = 64 << 10

func FuzzGoAMXLoadBytes(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("AMX"))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > maxFuzzImageSize {
			data = data[:maxFuzzImageSize]
		}

		vm, err := (GoAMXBackend{}).LoadBytes("fuzz.amx", data)
		if err == nil {
			_ = vm.Close()
		}
	})
}
