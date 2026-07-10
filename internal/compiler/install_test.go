package compiler

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInstallOpenMPCompilerIntegration(t *testing.T) {
	if os.Getenv("PAWNTEST_INSTALL_INTEGRATION") == "" {
		t.Skip("set PAWNTEST_INSTALL_INTEGRATION to test the GitHub release download")
	}

	path, err := InstallOpenMPCompiler(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if err := validateCompiler(path); err != nil {
		t.Fatal(err)
	}

	if path.Path == "" {
		t.Fatal("installed compiler has empty Path")
	}
}

func TestSelectAssetChoosesCurrentPlatformCompiler(t *testing.T) {
	assets := []releaseAsset{
		{Name: "pawnc-3.10.11-linux.tar.gz", BrowserDownloadURL: "linux"},
		{Name: "pawnc-3.10.11-windows.zip", BrowserDownloadURL: "windows"},
		{Name: "pawnc-3.10.11-macos.tar.gz", BrowserDownloadURL: "macos"},
	}

	tests := []struct {
		goos string
		want string
	}{
		{goos: "linux", want: "linux"},
		{goos: "windows", want: "windows"},
		{goos: "darwin", want: "macos"},
	}
	for _, tt := range tests {
		got, err := selectAsset(assets, tt.goos, "amd64")
		if err != nil {
			t.Fatal(err)
		}

		if got.BrowserDownloadURL != tt.want {
			t.Fatalf("selectAsset(%q) = %q, want %q", tt.goos, got.BrowserDownloadURL, tt.want)
		}
	}
}

func TestSelectAssetPrefersMatchingArch(t *testing.T) {
	assets := []releaseAsset{
		{Name: "pawnc-3.10.11-linux-x86.tar.gz", BrowserDownloadURL: "x86"},
		{Name: "pawnc-3.10.11-linux-x86_64.tar.gz", BrowserDownloadURL: "x64"},
	}

	got, err := selectAsset(assets, "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}

	if got.BrowserDownloadURL != "x64" {
		t.Fatalf("selectAsset() = %q, want x64", got.BrowserDownloadURL)
	}
}

func TestSelectAssetRejectsGenericX86AssetOnARM(t *testing.T) {
	assets := []releaseAsset{{Name: "pawnc-3.10.11-linux.tar.gz", BrowserDownloadURL: "x86"}}
	if _, err := selectAsset(assets, "linux", "arm64"); !errors.Is(err, ErrNoCompilerAsset) {
		t.Fatalf("selectAsset() error = %v, want ErrNoCompilerAsset", err)
	}
}

func TestSelectAssetDoesNotTreatDarwinAsWindows(t *testing.T) {
	assets := []releaseAsset{
		{Name: "pawnc-3.10.11-darwin.tar.gz", BrowserDownloadURL: "darwin"},
		{Name: "pawnc-3.10.11-windows.zip", BrowserDownloadURL: "windows"},
	}

	got, err := selectAsset(assets, "windows", "amd64")
	if err != nil {
		t.Fatal(err)
	}

	if got.BrowserDownloadURL != "windows" {
		t.Fatalf("selectAsset() = %q, want windows", got.BrowserDownloadURL)
	}
}

func TestFindPawnCCAndMakeCompiler(t *testing.T) {
	root := t.TempDir()

	bin := filepath.Join(root, "archive", "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}

	name := "pawncc"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	realPawnCC := filepath.Join(bin, name)
	if err := os.WriteFile(realPawnCC, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	found, err := findPawnCC(root)
	if err != nil {
		t.Fatal(err)
	}

	if found != realPawnCC {
		t.Fatalf("findPawnCC() = %q, want %q", found, realPawnCC)
	}

	c := makeCompiler(found)
	if c.Path != found {
		t.Fatalf("makeCompiler().Path = %q, want %q", c.Path, found)
	}

	wantLib := filepath.Join(root, "archive", "lib")
	if len(c.LibDirs) < 2 || c.LibDirs[1] != wantLib {
		t.Fatalf("makeCompiler().LibDirs = %#v, want sibling lib %q first", c.LibDirs, wantLib)
	}

	cmd := c.Command("--help")
	if cmd.Path != found {
		t.Fatalf("Command().Path = %q, want %q", cmd.Path, found)
	}

	var ldLine string

	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "LD_LIBRARY_PATH=") {
			ldLine = e
			break
		}
	}

	if !strings.Contains(ldLine, wantLib) {
		t.Fatalf("LD_LIBRARY_PATH env = %q, want it to include %q", ldLine, wantLib)
	}
}

func TestBareCompilerHasNoLibDirs(t *testing.T) {
	c := Bare("/usr/bin/pawncc")
	if c.Path != "/usr/bin/pawncc" || len(c.LibDirs) != 0 {
		t.Fatalf("Bare() = %+v, want empty LibDirs", c)
	}

	cmd := c.Command("a")
	if cmd.Path != "/usr/bin/pawncc" {
		t.Fatalf("Bare().Command().Path = %q", cmd.Path)
	}

	if cmd.Env != nil {
		t.Fatalf("Bare().Command().Env = %#v, want nil (inherit process env)", cmd.Env)
	}
}

func TestExtractTarGZRejectsUnsafePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.tar.gz")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	gz := gzip.NewWriter(f)

	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "../pawncc", Mode: 0o755, Size: 0, Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if err := extractTarGZ(path, t.TempDir()); err == nil {
		t.Fatal("expected unsafe archive path error")
	}
}

func TestExtractZipContinuesAfterDirectoryEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "compiler.zip")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	zw := zip.NewWriter(f)
	if _, err := zw.Create("pawnc/"); err != nil {
		t.Fatal(err)
	}

	w, err := zw.Create("pawnc/pawncc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := w.Write([]byte("binary")); err != nil {
		t.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(dir, "out")
	if err := extractZip(path, dest); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dest, "pawnc", "pawncc")); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyAssetDigest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "asset")

	data := []byte("compiler")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	sum := sha256.Sum256(data)
	if err := verifyAssetDigest(path, fmt.Sprintf("sha256:%x", sum[:])); err != nil {
		t.Fatal(err)
	}

	if err := verifyAssetDigest(path, "sha256:0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatal("expected checksum mismatch")
	}

	if err := verifyAssetDigest(path, ""); !errors.Is(err, ErrMissingCompilerDigest) {
		t.Fatalf("missing digest error = %v, want ErrMissingCompilerDigest", err)
	}
}

func TestCompilerValidationDetailsIncludesOutputAndLinuxHint(t *testing.T) {
	got := compilerValidationDetails("missing /lib/ld-linux.so.2", "linux")
	for _, want := range []string{
		"missing /lib/ld-linux.so.2",
		"32-bit glibc binary",
		"lib32-glibc",
		"libc6-i386",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("compilerValidationDetails() missing %q:\n%s", want, got)
		}
	}
}

func TestCompilerStartedAcceptsPawnUsageBanner(t *testing.T) {
	output := "Pawn compiler 3.10.11 Copyright (c) 1997-2006, ITB CompuPhase\n\nUsage:   pawncc <filename> [filename...] [options]\n"
	if !compilerStarted(output) {
		t.Fatal("compilerStarted rejected pawncc usage banner")
	}

	if compilerStarted("/lib/ld-linux.so.2: bad ELF interpreter") {
		t.Fatal("compilerStarted accepted dynamic loader error")
	}
}
