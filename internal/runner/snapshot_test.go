package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotStoreWritesAndReloads(t *testing.T) {
	source := filepath.Join(t.TempDir(), "value.test.pwn")
	store := newSnapshotStore(source, true)
	store.values["test_value::output"] = "first"
	if err := store.write(); err != nil {
		t.Fatal(err)
	}
	store.values["test_value::output"] = "second"
	if err := store.write(); err != nil {
		t.Fatal(err)
	}
	reloaded := newSnapshotStore(source, false)
	if reloaded.loadErr != nil {
		t.Fatal(reloaded.loadErr)
	}
	if got := reloaded.values["test_value::output"]; got != "second" {
		t.Fatalf("snapshot = %q, want second", got)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(source), "__snapshots__", "value.test.pwn.snap.json")); err != nil {
		t.Fatal(err)
	}
}
