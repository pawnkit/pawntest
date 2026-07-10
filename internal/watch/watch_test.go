package watch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilesDetectsRelevantContentChanges(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "math.test.pwn")
	ignored := filepath.Join(dir, "notes.txt")

	if err := os.WriteFile(source, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(ignored, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}

	before := Files([]string{dir})

	if err := os.WriteFile(ignored, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !Equal(before, Files([]string{dir})) {
		t.Fatal("irrelevant file changed snapshot")
	}

	if err := os.WriteFile(source, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}

	if Equal(before, Files([]string{dir})) {
		t.Fatal("Pawn source change was not detected")
	}
}

func TestPathsNormalizesRecursiveRootsAndWatchesSourceDirectory(t *testing.T) {
	dir := t.TempDir()

	source := filepath.Join(dir, "math.test.pwn")
	if err := os.WriteFile(source, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	paths := Paths([]string{source, dir + "/..."}, nil)
	foundDir := false

	for _, path := range paths {
		if path == dir {
			foundDir = true
		}

		if filepath.Base(path) == "..." {
			t.Fatalf("recursive suffix was not normalized: %q", path)
		}
	}

	if !foundDir {
		t.Fatalf("source directory missing from paths: %#v", paths)
	}
}
