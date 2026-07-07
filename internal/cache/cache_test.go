package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDirUsesPawntestCacheSubdirectory(t *testing.T) {
	dir := filepath.Clean(Dir())
	if filepath.Base(dir) != "pawntest" && !strings.HasSuffix(dir, "pawntest-cache") {
		t.Fatalf("Dir() = %q, want pawntest cache directory", dir)
	}
}

func TestIncludeDirInWritesEmbeddedIncludeWhenMissing(t *testing.T) {
	dir, err := IncludeDirIn(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "pawntest.inc"))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("cached pawntest.inc is empty")
	}
	for _, relative := range []string{"pawntest/core.inc", "pawntest/assertions.inc", "pawntest/mocks.inc", "pawntest/scenarios.inc"} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(relative))); err != nil {
			t.Fatalf("embedded include module %s was not extracted: %v", relative, err)
		}
	}
}

func TestIncludeDirInDoesNotRewriteUnchangedInclude(t *testing.T) {
	base := t.TempDir()
	dir, err := IncludeDirIn(base)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "pawntest.inc")
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if _, err := IncludeDirIn(base); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("IncludeDirIn rewrote unchanged include: before %s after %s", before.ModTime(), after.ModTime())
	}
}

func TestIncludeDirInRefreshesStaleInclude(t *testing.T) {
	base := t.TempDir()
	dir, err := IncludeDirIn(base)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "pawntest.inc")
	if err := os.WriteFile(path, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := IncludeDirIn(base); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "stale" {
		t.Fatal("IncludeDirIn did not refresh stale include")
	}
}

func TestEmbeddedNativeNamesFitPawnSymbolLimit(t *testing.T) {
	data := IncludeBundleBytes()
	for lineNumber, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "native ") {
			continue
		}
		name := strings.TrimPrefix(line, "native ")
		if index := strings.IndexByte(name, '('); index >= 0 {
			name = name[:index]
		}
		if len(name) > 31 {
			t.Errorf("pawntest.inc:%d native %q exceeds Pawn's 31-character symbol limit", lineNumber+1, name)
		}
	}
}

func TestEmbeddedIncludeDoesNotExposeRemovedMockAPI(t *testing.T) {
	include := string(IncludeBundleBytes())
	for _, removed := range []string{
		"MOCK_CALL_COUNT", "MOCK_ARG", "EXPECT_CALLS(", "EXPECT_ARG(",
		"EXPECT_STRING_ARG", "EXPECT_CALL_ORDER", "__pawntest_mock_call_count", "__pawntest_mock_arg",
	} {
		if strings.Contains(include, removed) {
			t.Errorf("embedded include still exposes removed API %q", removed)
		}
	}
}
