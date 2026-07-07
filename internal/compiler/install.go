package compiler

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const OpenMPCompilerReleaseAPI = "https://api.github.com/repos/openmultiplayer/compiler/releases/latest"
const maxCompilerDownloadBytes = 200 << 20

var ErrNoCompilerAsset = errors.New("no compatible openmultiplayer compiler asset found")
var ErrMissingCompilerDigest = errors.New("compiler release asset has no SHA-256 digest")

type releaseInfo struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
}

func InstallOpenMPCompiler(ctx context.Context, dir string) (*Compiler, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	release, err := fetchRelease(ctx, client, OpenMPCompilerReleaseAPI)
	if err != nil {
		return nil, err
	}
	asset, err := selectAsset(release.Assets, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	root := filepath.Join(dir, "tools", "openmp-compiler", strings.TrimPrefix(release.TagName, "v"), runtime.GOOS)
	if cached, ok := loadCachedCompiler(root); ok {
		if err := validateCompiler(cached); err == nil {
			return cached, nil
		}
		if realPawnCC, err := findPawnCC(root); err == nil {
			repaired := makeCompiler(realPawnCC)
			if err := writeCompilerMarker(root, realPawnCC); err != nil {
				return nil, err
			}
			return repaired, validateCompiler(repaired)
		}
	}
	parent := filepath.Dir(root)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return nil, err
	}
	tmp, err := downloadAsset(ctx, client, asset.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp)
	stage, err := os.MkdirTemp(parent, ".download-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(stage)
	if err := verifyAssetDigest(tmp, asset.Digest); err != nil {
		return nil, err
	}
	if err := extractAsset(tmp, asset.Name, stage); err != nil {
		return nil, err
	}
	if _, err := findPawnCC(stage); err != nil {
		return nil, err
	}
	if err := os.RemoveAll(root); err != nil {
		return nil, err
	}
	if err := os.Rename(stage, root); err != nil {
		return nil, err
	}
	realPawnCC, err := findPawnCC(root)
	if err != nil {
		return nil, err
	}
	if err := writeCompilerMarker(root, realPawnCC); err != nil {
		return nil, err
	}
	c := makeCompiler(realPawnCC)
	return c, validateCompiler(c)
}
