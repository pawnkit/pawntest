package cache

import (
	"bytes"
	"os"
	"path/filepath"

	pawntestinclude "github.com/pawnkit/pawntest/include"
)

func IncludeBytes() ([]byte, error) {
	return pawntestinclude.Pawntest(), nil
}

func IncludeBundleBytes() []byte {
	files := pawntestinclude.Files()
	paths := []string{"pawntest.inc", "pawntest/core.inc", "pawntest/assertions.inc", "pawntest/mocks.inc", "pawntest/scenarios.inc"}
	var bundle []byte
	for _, path := range paths {
		bundle = append(bundle, []byte(path)...)
		bundle = append(bundle, 0)
		bundle = append(bundle, files[path]...)
		bundle = append(bundle, 0)
	}
	return bundle
}

func Dir() string {
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		return filepath.Join(os.TempDir(), "pawntest-cache")
	}
	return filepath.Join(dir, "pawntest")
}

func IncludeDir() (string, error) {
	return IncludeDirIn(Dir())
}

func IncludeDirIn(base string) (string, error) {
	if base == "" {
		base = Dir()
	}
	dir := filepath.Join(base, "include")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for path, data := range pawntestinclude.Files() {
		destination := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return "", err
		}
		if err := writeIfMissingOrChanged(destination, data, 0o644); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func writeIfMissingOrChanged(path string, data []byte, perm os.FileMode) error {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, data) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func AMXDirIn(base string) (string, error) {
	if base == "" {
		base = Dir()
	}
	dir := filepath.Join(base, "amx")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
