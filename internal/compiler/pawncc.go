package compiler

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pawnkit/pawntest/internal/cache"
)

func Compile(path string, opts Options) (string, error) {
	return CompileContext(context.Background(), path, opts)
}

func CompileContext(ctx context.Context, path string, opts Options) (string, error) {
	c := opts.Compiler

	pawnccPath := "pawncc"
	if c != nil && c.Path != "" {
		pawnccPath = c.Path
	}

	outDir := opts.OutDir
	if outDir == "" {
		outDir = filepath.Dir(path)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}

	key, err := compileKey(path, pawnccPath, opts)
	if err != nil {
		return "", err
	}

	outName := key + ".amx"
	if opts.NoCache {
		outName = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".amx"
	}

	out := filepath.Join(outDir, outName)
	if !opts.NoCache && opts.Count != 1 {
		if _, err := os.Stat(out); err == nil {
			return out, nil
		}
	}

	args := buildArgs(out, path, opts)

	var cmd *exec.Cmd
	if c != nil {
		cmd = c.CommandContext(ctx, args...)
	} else {
		cmd = exec.CommandContext(ctx, pawnccPath, args...)
	}

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if _, lookErr := exec.LookPath(pawnccPath); lookErr != nil && pawnccPath == "pawncc" {
			return "", ErrPawnCCNotFound
		}

		return "", fmt.Errorf("pawncc failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return out, nil
}

func buildArgs(out, path string, opts Options) []string {
	args := []string{"-C-", "-d2", "-o" + out}
	for _, inc := range opts.Includes {
		args = append(args, "-i"+inc)
	}

	for _, def := range opts.Defines {
		args = append(args, "-D"+def)
	}

	args = append(args, opts.ExtraArgs...)
	args = append(args, path)

	return args
}

func compileKey(path, pawncc string, opts Options) (string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var parts [][]byte

	parts = append(parts, src)
	parts = append(parts, []byte(pawncc))
	parts = append(parts, []byte("pawntest-compiler-mode=-C-,-d2"))

	parts = append(parts, cache.IncludeBundleBytes())
	for _, inc := range opts.Includes {
		parts = append(parts, []byte("include="+inc))
	}

	includeDeps, err := includeDependencyBytes(path, opts.Includes)
	if err != nil {
		return "", err
	}

	parts = append(parts, includeDeps...)
	for _, def := range opts.Defines {
		parts = append(parts, []byte("define="+def))
	}

	for _, arg := range opts.ExtraArgs {
		parts = append(parts, []byte("arg="+arg))
	}

	return cache.Key(parts...), nil
}

func includeDependencyBytes(src string, includeDirs []string) ([][]byte, error) {
	seen := map[string]bool{}

	search := append([]string{filepath.Dir(src)}, includeDirs...)

	return collectIncludeDependencies(src, search, seen)
}

func collectIncludeDependencies(path string, search []string, seen map[string]bool) ([][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parts [][]byte

	for _, include := range parseIncludes(string(data)) {
		includePath, ok := resolveInclude(include, search)
		if !ok || seen[includePath] {
			continue
		}

		seen[includePath] = true

		includeData, err := os.ReadFile(includePath)
		if err != nil {
			return nil, err
		}

		parts = append(parts, []byte("include-file="+includePath), includeData)

		nested, err := collectIncludeDependencies(includePath, search, seen)
		if err != nil {
			return nil, err
		}

		parts = append(parts, nested...)
	}

	return parts, nil
}

func parseIncludes(src string) []string {
	var out []string

	for line := range strings.SplitSeq(src, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#include") {
			continue
		}

		rest := strings.TrimSpace(strings.TrimPrefix(line, "#include"))
		if rest == "" {
			continue
		}

		if rest[0] == '<' {
			if end := strings.IndexByte(rest, '>'); end > 1 {
				out = append(out, rest[1:end])
			}

			continue
		}

		if rest[0] == '"' {
			if end := strings.IndexByte(rest[1:], '"'); end >= 0 {
				out = append(out, rest[1:1+end])
			}
		}
	}

	return out
}

func resolveInclude(name string, search []string) (string, bool) {
	candidates := []string{name}
	if filepath.Ext(name) == "" {
		candidates = append(candidates, name+".inc")
	}

	for _, dir := range search {
		for _, candidate := range candidates {
			path := filepath.Join(dir, candidate)
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				abs, err := filepath.Abs(path)
				if err == nil {
					return abs, true
				}

				return path, true
			}
		}
	}

	return "", false
}
