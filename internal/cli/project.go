package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/pawnkit/pawn-project/fsx"
	"github.com/pawnkit/pawn-project/manifest"
	"github.com/pawnkit/pawn-project/workspace"
	"github.com/pawnkit/pawnkit-core/source"
)

func loadProjectConfig(start string) (Config, string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return Config{}, "", err
	}

	root, err := workspace.FindRoot(fsx.OS{}, abs)
	if err != nil {
		if errors.Is(err, workspace.ErrNotFound) {
			return Config{}, "", nil
		}

		return Config{}, "", err
	}

	result, err := manifest.Load(source.NewRegistry(), fsx.OS{}, root.ManifestPath)
	if err != nil {
		return Config{}, "", err
	}

	if result.Manifest == nil {
		return Config{}, "", nil
	}

	rootDir := filepath.Clean(root.Dir)
	cfg := Config{Include: absolutePaths(rootDir, result.Manifest.EffectiveIncludePaths())}

	if err := applyProjectExtension(&cfg, result.Manifest.PawnKit); err != nil {
		return Config{}, "", err
	}

	cfg.Include = absolutePaths(rootDir, cfg.Include)
	cfg.Tests = absolutePaths(rootDir, cfg.Tests)

	return cfg, rootDir, cfg.validate()
}

func applyProjectExtension(cfg *Config, extension *manifest.PawnKitExtension) error {
	if extension == nil {
		return nil
	}

	if tool := extension.Tool["pawntest"]; tool != nil {
		data, err := json.Marshal(tool)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(cfg); err != nil {
			return err
		}
	}

	if paths, ok := stringSlice(extension.Tests["paths"]); ok && len(cfg.Tests) == 0 {
		cfg.Tests = paths
	}

	return nil
}

func absolutePaths(root string, paths []string) []string {
	out := make([]string, len(paths))
	for i, path := range paths {
		if filepath.IsAbs(path) {
			out[i] = filepath.Clean(path)
		} else {
			out[i] = filepath.Join(root, path)
		}
	}

	return out
}

func stringSlice(value any) ([]string, bool) {
	items, ok := value.([]any)
	if !ok {
		return nil, false
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok {
			return nil, false
		}

		out = append(out, text)
	}

	return out, true
}

func workingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return wd
}
