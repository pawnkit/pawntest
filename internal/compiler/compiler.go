package compiler

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

// Compiler configures pawncc and its library paths.
type Compiler struct {
	Path    string
	LibDirs []string
}

// Bare configures pawncc without inferred library paths.
func Bare(path string) *Compiler {
	return &Compiler{Path: path}
}

// FromPath configures pawncc with adjacent and sibling library directories.
func FromPath(path string) *Compiler {
	return &Compiler{Path: path, LibDirs: compilerLibraryDirs(path)}
}

// Command creates a pawncc command.
func (c *Compiler) Command(args ...string) *exec.Cmd {
	return c.CommandContext(context.Background(), args...)
}

// CommandContext creates a cancellable pawncc command.
func (c *Compiler) CommandContext(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.Path, args...)
	if len(c.LibDirs) == 0 {
		return cmd
	}

	list := strings.Join(c.LibDirs, string(os.PathListSeparator))

	env := append([]string{}, os.Environ()...)
	env = appendLibraryPath(env, "LD_LIBRARY_PATH", list)
	env = appendLibraryPath(env, "DYLD_LIBRARY_PATH", list)
	cmd.Env = env

	return cmd
}

// appendLibraryPath adds a directory to a library path variable.
func appendLibraryPath(env []string, name, value string) []string {
	prefix := name + "="
	for i, e := range env {
		if after, ok := strings.CutPrefix(e, prefix); ok {
			existing := after
			if existing != "" {
				value = value + string(os.PathListSeparator) + existing
			}

			env[i] = prefix + value

			return env
		}
	}

	return append(env, prefix+value)
}
