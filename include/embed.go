package include

import (
	"embed"
	"io/fs"
)

//go:embed pawntest.inc pawntest/*.inc pawntest/scenarios/*.inc
var files embed.FS

func Pawntest() []byte {
	data, _ := files.ReadFile("pawntest.inc")
	return append([]byte(nil), data...)
}

func Files() map[string][]byte {
	out := map[string][]byte{}

	_ = fs.WalkDir(files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		data, _ := files.ReadFile(path)
		out[path] = append([]byte(nil), data...)

		return nil
	})

	return out
}
