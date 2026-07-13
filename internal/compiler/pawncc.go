package compiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pawnkit/pawntest/internal/cache"
)

func Compile(path string, opts Options) (string, error) {
	return CompileContext(context.Background(), path, opts)
}

func CompileContext(ctx context.Context, path string, opts Options) (string, error) {
	if err := validateGeneratedSymbols(path); err != nil {
		return "", err
	}

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

	unlock, err := lockCompile(ctx, out+".lock")
	if err != nil {
		return "", err
	}
	defer unlock()

	if !opts.NoCache && opts.Count != 1 {
		if _, err := os.Stat(out); err == nil {
			return out, nil
		}
	}

	temp, err := os.CreateTemp(outDir, ".pawntest-*.amx")
	if err != nil {
		return "", err
	}

	tempOut := temp.Name()
	if err := temp.Close(); err != nil {
		os.Remove(tempOut)
		return "", err
	}

	if err := os.Remove(tempOut); err != nil {
		return "", err
	}
	defer os.Remove(tempOut)

	args := buildArgs(tempOut, path, opts)

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

	if opts.Diagnostics != nil && stderr.Len() > 0 {
		if err := writeCompilerDiagnostics(opts.Diagnostics, stderr.String()); err != nil {
			return "", err
		}
	}

	if err := replaceFile(tempOut, out); err != nil {
		return "", err
	}

	return out, nil
}

func writeCompilerDiagnostics(writer io.Writer, diagnostics string) error {
	for line := range strings.SplitSeq(diagnostics, "\n") {
		if strings.Contains(line, `warning 209: function "test_`) {
			continue
		}

		if line == "" {
			continue
		}

		if _, err := fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}

	return nil
}

func lockCompile(ctx context.Context, path string) (func(), error) {
	for {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			_ = file.Close()
			return func() { _ = os.Remove(path) }, nil
		}

		if !os.IsExist(err) {
			return nil, err
		}

		if info, statErr := os.Stat(path); statErr == nil && time.Since(info.ModTime()) > 30*time.Minute {
			_ = os.Remove(path)
			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(25 * time.Millisecond):
		}
	}
}

func replaceFile(source, destination string) error {
	if err := os.Rename(source, destination); err == nil {
		return nil
	}

	if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.Rename(source, destination)
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
	parts = append(parts, compilerIdentity(pawncc, opts.Compiler))
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

func compilerIdentity(path string, configured *Compiler) []byte {
	resolved := path
	if found, err := exec.LookPath(path); err == nil {
		resolved = found
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return []byte("compiler-path=" + path)
	}

	identity := append([]byte("compiler="+resolved+"\x00"), data...)
	if configured == nil {
		return identity
	}

	for _, dir := range configured.LibDirs {
		for _, name := range []string{"libpawnc.so", "libpawnc.dylib", "pawnc.dll", "libpawnc.dll"} {
			library := filepath.Join(dir, name)

			data, err := os.ReadFile(library)
			if err == nil {
				identity = append(identity, []byte("\x00library="+library+"\x00")...)
				identity = append(identity, data...)
			}
		}
	}

	return identity
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

	for _, include := range parseIncludeDirectives(string(data)) {
		includeSearch := search
		if include.quoted {
			includeSearch = append([]string{filepath.Dir(path)}, search...)
		}

		includePath, ok := resolveInclude(include.name, includeSearch)
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
	for _, include := range parseIncludeDirectives(src) {
		out = append(out, include.name)
	}

	return out
}

type includeDirective struct {
	name   string
	quoted bool
}

func parseIncludeDirectives(src string) []includeDirective {
	var out []includeDirective

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
				out = append(out, includeDirective{name: rest[1:end]})
			}

			continue
		}

		if rest[0] == '"' {
			if end := strings.IndexByte(rest[1:], '"'); end >= 0 {
				out = append(out, includeDirective{name: rest[1 : 1+end], quoted: true})
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
