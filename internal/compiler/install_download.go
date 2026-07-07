package compiler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadAsset(ctx context.Context, client *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("compiler download failed: %s", resp.Status)
	}

	tmp, err := os.CreateTemp("", "pawntest-compiler-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, io.LimitReader(resp.Body, maxCompilerDownloadBytes+1)); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	if info, err := tmp.Stat(); err != nil {
		os.Remove(tmp.Name())
		return "", err
	} else if info.Size() > maxCompilerDownloadBytes {
		os.Remove(tmp.Name())
		return "", errors.New("compiler download exceeded maximum allowed size")
	}
	return tmp.Name(), nil
}

func verifyAssetDigest(path, digest string) error {
	if digest == "" {
		return ErrMissingCompilerDigest
	}

	alg, value, ok := strings.Cut(digest, ":")
	if !ok || strings.ToLower(alg) != "sha256" {
		return fmt.Errorf("%w: %q", ErrMissingCompilerDigest, digest)
	}

	want, err := hex.DecodeString(value)
	if err != nil {
		return fmt.Errorf("invalid release asset digest %q: %w", digest, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return err
	}
	if !bytes.Equal(hash.Sum(nil), want) {
		return errors.New("compiler download checksum mismatch")
	}
	return nil
}
