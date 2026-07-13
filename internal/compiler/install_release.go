package compiler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	platformWindows = "windows"
	platformDarwin  = "darwin"
	platformLinux   = "linux"
)

func fetchRelease(ctx context.Context, client *http.Client, url string) (releaseInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return releaseInfo{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return releaseInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return releaseInfo{}, fmt.Errorf("GitHub release request failed: %s", resp.Status)
	}

	var release releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return releaseInfo{}, err
	}

	return release, nil
}

func selectAsset(assets []releaseAsset, goos, goarch string) (releaseAsset, error) {
	needles, err := platformAssetNeedles(goos)
	if err != nil {
		return releaseAsset{}, err
	}

	var fallback *releaseAsset

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if !isCompilerAsset(name) {
			continue
		}

		for _, needle := range needles {
			if !strings.Contains(name, needle) {
				continue
			}

			if assetMatchesArch(name, goarch) {
				return asset, nil
			}

			if fallback == nil && supportsOpenMPX86Fallback(goarch) {
				candidate := asset
				fallback = &candidate
			}
		}
	}

	if fallback != nil {
		return *fallback, nil
	}

	return releaseAsset{}, fmt.Errorf("%s: %w", goos, ErrNoCompilerAsset)
}

func platformAssetNeedles(goos string) ([]string, error) {
	switch goos {
	case platformWindows:
		return []string{platformWindows, "win32", "win64"}, nil
	case platformDarwin:
		return []string{"mac", "macos", platformDarwin}, nil
	case platformLinux:
		return []string{platformLinux}, nil
	default:
		return nil, fmt.Errorf("%s: %w", goos, ErrNoCompilerAsset)
	}
}

func isCompilerAsset(name string) bool {
	return strings.Contains(name, "pawnc") && !strings.Contains(name, "source")
}

func assetMatchesArch(name, goarch string) bool {
	archTerms := map[string][]string{
		"amd64": {"amd64", "x86_64", "x64"},
		"386":   {"386", "i386", "x86"},
		"arm64": {"arm64", "aarch64"},
	}

	terms := archTerms[goarch]
	if len(terms) == 0 {
		return true
	}

	if !assetNamesAnyKnownArch(name, archTerms) {
		return supportsOpenMPX86Fallback(goarch)
	}

	for _, term := range terms {
		if strings.Contains(name, term) {
			return true
		}
	}

	return false
}

func assetNamesAnyKnownArch(name string, archTerms map[string][]string) bool {
	for _, terms := range archTerms {
		for _, term := range terms {
			if strings.Contains(name, term) {
				return true
			}
		}
	}

	return false
}

func supportsOpenMPX86Fallback(goarch string) bool {
	return goarch == "amd64" || goarch == "386"
}
