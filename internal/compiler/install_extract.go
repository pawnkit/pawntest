package compiler

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractAsset(path, name, dest string) error {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		return extractZip(path, dest)
	case strings.HasSuffix(lower, ".tar.gz"), strings.HasSuffix(lower, ".tgz"):
		return extractTarGZ(path, dest)
	default:
		return fmt.Errorf("unsupported compiler archive %q", name)
	}
}

func extractZip(path, dest string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := safeExtractPath(dest, f.Name); err != nil {
			return err
		}

		target := filepath.Join(dest, filepath.Clean(f.Name))
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			return err
		}
		if err := writeFileFromReader(target, src, f.FileInfo().Mode()); err != nil {
			src.Close()
			return err
		}
		src.Close()
	}
	return nil
}

func extractTarGZ(path, dest string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := extractTarEntry(dest, hdr, tr); err != nil {
			return err
		}
	}
}

func extractTarEntry(dest string, hdr *tar.Header, r io.Reader) error {
	if err := safeExtractPath(dest, hdr.Name); err != nil {
		return err
	}

	target := filepath.Join(dest, filepath.Clean(hdr.Name))
	switch hdr.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(target, 0o755)
	case tar.TypeReg:
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return writeFileFromReader(target, r, os.FileMode(hdr.Mode))
	default:
		return nil
	}
}

func safeExtractPath(dest, name string) error {
	clean := filepath.Clean(name)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("unsafe archive path %q", name)
	}

	target := filepath.Join(dest, clean)
	rel, err := filepath.Rel(dest, target)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("unsafe archive path %q", name)
	}
	return nil
}

func writeFileFromReader(path string, r io.Reader, mode os.FileMode) error {
	if mode == 0 {
		mode = 0o644
	}

	out, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}
