package include

import "embed"

//go:embed pawntest.inc pawntest/*.inc
var files embed.FS

func Pawntest() []byte {
	data, _ := files.ReadFile("pawntest.inc")
	return append([]byte(nil), data...)
}

func Files() map[string][]byte {
	out := map[string][]byte{}
	for _, path := range []string{"pawntest.inc", "pawntest/core.inc", "pawntest/assertions.inc", "pawntest/mocks.inc", "pawntest/scenarios.inc"} {
		data, _ := files.ReadFile(path)
		out[path] = append([]byte(nil), data...)
	}
	return out
}
