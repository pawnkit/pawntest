package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Files(patterns []string) ([]string, error) {
	if len(patterns) == 0 {
		patterns = []string{"."}
	}
	seen := map[string]bool{}
	var out []string
	for _, p := range patterns {
		matches, err := expand(p)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	return out, nil
}

func expand(path string) ([]string, error) {
	recursive := strings.HasSuffix(path, "/...")
	root := strings.TrimSuffix(path, "/...")
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if isPawnTest(path) || strings.EqualFold(filepath.Ext(path), ".amx") {
			return []string{path}, nil
		}
		return nil, nil
	}
	var out []string
	walk := func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !recursive && p != root {
				return filepath.SkipDir
			}
			return nil
		}
		if isPawnTest(p) || strings.EqualFold(filepath.Ext(p), ".amx") {
			out = append(out, p)
		}
		return nil
	}
	return out, filepath.WalkDir(root, walk)
}

func isPawnTest(path string) bool {
	base := filepath.Base(path)
	lower := strings.ToLower(base)
	return strings.HasSuffix(lower, ".test.pwn") || strings.HasSuffix(lower, ".test.inc")
}
