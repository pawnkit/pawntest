package compiler

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func makeCompiler(pawnccPath string) *Compiler {
	return FromPath(pawnccPath)
}

func FindCachedCompiler(cacheDir string) (*Compiler, bool) {
	root := filepath.Join(cacheDir, "tools", "openmp-compiler")

	var found *Compiler

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || found != nil || d.IsDir() || filepath.Base(path) != "pawncc.path" {
			return nil
		}

		compiler, ok := compilerFromMarker(path)
		if !ok {
			return nil
		}

		found = compiler

		return filepath.SkipAll
	})

	return found, found != nil
}

func loadCachedCompiler(root string) (*Compiler, bool) {
	return compilerFromMarker(markerPath(root))
}

func compilerFromMarker(path string) (*Compiler, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	realPawnCC := strings.TrimSpace(string(data))
	if realPawnCC == "" {
		return nil, false
	}

	if _, err := os.Stat(realPawnCC); err != nil {
		return nil, false
	}

	return makeCompiler(realPawnCC), true
}

func writeCompilerMarker(root, pawnccPath string) error {
	dir := pawnccBinDir(root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(markerPath(root), []byte(pawnccPath+"\n"), 0o644)
}

func findPawnCC(root string) (string, error) {
	name := defaultPawnCC
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	var found string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if !strings.EqualFold(filepath.Base(path), name) {
			return nil
		}

		if filepath.Base(filepath.Dir(path)) == "pawncc-bin" {
			return nil
		}

		found = path

		return filepath.SkipAll
	})
	if err != nil {
		return "", err
	}

	if found == "" {
		return "", errors.New("downloaded compiler archive did not contain pawncc")
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(found, 0o755); err != nil {
			return "", err
		}
	}

	return found, nil
}

func pawnccBinDir(root string) string {
	return filepath.Join(root, "pawncc-bin")
}

func markerPath(root string) string {
	return filepath.Join(pawnccBinDir(root), "pawncc.path")
}

func compilerLibraryDirs(pawncc string) []string {
	binDir := filepath.Dir(pawncc)
	dirs := []string{binDir}

	siblingLib := filepath.Join(filepath.Dir(binDir), "lib")
	if siblingLib != binDir {
		dirs = append(dirs, siblingLib)
	}

	return dirs
}
