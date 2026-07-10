package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilesDiscoversPawnTestsAndAMX(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "math.test.pwn"))
	mustWrite(t, filepath.Join(dir, "helpers.test.inc"))
	mustWrite(t, filepath.Join(dir, "legacy_test.pwn"))
	mustWrite(t, filepath.Join(dir, "test_legacy.pwn"))
	mustWrite(t, filepath.Join(dir, "skip.pwn"))
	mustWrite(t, filepath.Join(dir, "compiled.amx"))

	nested := filepath.Join(dir, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	mustWrite(t, filepath.Join(nested, "deep.test.pwn"))

	got, err := Files([]string{dir})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{filepath.Join(dir, "compiled.amx"), filepath.Join(dir, "helpers.test.inc"), filepath.Join(dir, "math.test.pwn")}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("non-recursive mismatch (-want +got):\n%s", diff)
	}

	got, err = Files([]string{dir + "/..."})
	if err != nil {
		t.Fatal(err)
	}

	want = []string{filepath.Join(dir, "compiled.amx"), filepath.Join(dir, "helpers.test.inc"), filepath.Join(dir, "math.test.pwn"), filepath.Join(nested, "deep.test.pwn")}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("recursive mismatch (-want +got):\n%s", diff)
	}
}

func mustWrite(t *testing.T, path string) {
	t.Helper()

	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
}
