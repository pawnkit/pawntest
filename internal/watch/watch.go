package watch

import (
	"context"
	"crypto/sha256"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Snapshot map[string][sha256.Size]byte

func Files(paths []string) Snapshot {
	files := map[string][sha256.Size]byte{}
	seen := map[string]bool{}
	for _, root := range paths {
		if root == "" || seen[root] {
			continue
		}
		seen[root] = true
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			addFile(files, root)
			continue
		}
		_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if entry.IsDir() {
				if path != root && ignoredDirectory(entry.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			if watchedFile(path) {
				addFile(files, path)
			}
			return nil
		})
	}
	return files
}

func Wait(ctx context.Context, paths []string, previous Snapshot, interval time.Duration) (Snapshot, error) {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			next := Files(paths)
			if !Equal(previous, next) {
				return next, nil
			}
		}
	}
}

func Equal(left, right Snapshot) bool {
	if len(left) != len(right) {
		return false
	}
	for path, digest := range left {
		if right[path] != digest {
			return false
		}
	}
	return true
}

func Paths(inputs, includes []string) []string {
	var paths []string
	for _, input := range inputs {
		root := strings.TrimSuffix(input, string(filepath.Separator)+"...")
		root = strings.TrimSuffix(root, "/...")
		paths = append(paths, root)
		if info, err := os.Stat(root); err == nil && !info.IsDir() {
			paths = append(paths, filepath.Dir(root))
		}
	}
	paths = append(paths, includes...)
	paths = append(paths, "pawntest.json", "pawntest.yaml", "pawntest.yml", "pawn.json", "pawn.yaml", "pawn.yml")
	sort.Strings(paths)
	return paths
}

func addFile(files Snapshot, path string) {
	data, err := os.ReadFile(path)
	if err == nil {
		files[path] = sha256.Sum256(data)
	}
}

func watchedFile(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	return extension == ".pwn" || extension == ".inc" || extension == ".json" || extension == ".yaml" || extension == ".yml"
}

func ignoredDirectory(name string) bool {
	switch name {
	case ".git", ".hg", ".svn", "node_modules", "vendor", "build", "dist", "__snapshots__":
		return true
	default:
		return false
	}
}
