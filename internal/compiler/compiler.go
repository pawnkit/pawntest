package compiler

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

// Compiler describes how to invoke the pawncc compiler.
//
// Path is the real pawncc executable. LibDirs, when non-empty, are
// directories containing shared libraries the binary needs at load time
// (e.g. the 32-bit openmultiplayer Linux build). They are exposed to the
// process via LD_LIBRARY_PATH / DYLD_LIBRARY_PATH rather than by writing a
// shell wrapper, so the binary is launched directly with os/exec and the
// environment is controlled per-invocation.
type Compiler struct {
	Path    string
	LibDirs []string
}

// Bare returns a Compiler that invokes the binary at path with no extra
// library configuration. It is the zero-configuration form used when the
// caller supplies an explicit pawncc path (from the CLI, PATH lookup, etc.).
func Bare(path string) *Compiler {
	return &Compiler{Path: path}
}

// Command returns an *exec.Cmd configured to run pawncc with the given
// arguments. When LibDirs is non-empty the platform library-path environment
// variables are augmented so the binary's own shared libraries resolve.
func (c *Compiler) Command(args ...string) *exec.Cmd {
	return c.CommandContext(context.Background(), args...)
}

// CommandContext is like Command but honours ctx cancellation. The context is
// attached to the underlying exec.Cmd so that a cancelled context kills the
// pawncc process, allowing callers to interrupt long compilations.
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

// appendLibraryPath appends (or sets) the named library-path variable on env,
// preserving any pre-existing value by joining it with the OS path-list
// separator. This mirrors the shell idiom `${VAR:+:$VAR}`.
func appendLibraryPath(env []string, name, value string) []string {
	prefix := name + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			existing := strings.TrimPrefix(e, prefix)
			if existing != "" {
				value = value + string(os.PathListSeparator) + existing
			}
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}
