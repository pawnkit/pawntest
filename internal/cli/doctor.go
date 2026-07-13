package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
	"github.com/pawnkit/pawntest/internal/cache"
	"github.com/pawnkit/pawntest/internal/compiler"
	"github.com/pawnkit/pawntest/internal/runner"
)

func (d DoctorCmd) execute(_ context.Context, w io.Writer) error {
	cfg, cfgPath, err := LoadDefaultConfigWithPath()
	if err != nil {
		return err
	}

	explicitPawnCC := d.PawnCC != ""
	d.applyConfig(cfg)

	return d.doctor(w, cfgPath, explicitPawnCC)
}

func (d DoctorCmd) doctor(w io.Writer, configPath string, explicitPawnCC bool) error {
	cacheDir := d.CacheDir
	if cacheDir == "" {
		cacheDir = cache.Dir()
	}

	includeDir, includeErr := cache.IncludeDirIn(cacheDir)
	includeHash := "unavailable"

	if data, err := cache.IncludeBytes(); err == nil {
		sum := sha256.Sum256(data)
		includeHash = hex.EncodeToString(sum[:8])
	}

	pawncc, source, pawnccErr := d.resolveDoctorPawnCC(cacheDir, explicitPawnCC)

	version := "unavailable"
	versionOK := false
	if pawnccErr == nil {
		version, versionOK = compilerVersion(pawncc)
	}

	fmt.Fprintln(w, cliColor(w, "pawntest doctor", ansiBold))
	fmt.Fprintf(w, "platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	if configPath == "" {
		fmt.Fprintln(w, "config: none")
	} else {
		fmt.Fprintf(w, "config: %s\n", configPath)
	}

	fmt.Fprintf(w, "cache: %s\n", cacheDir)

	if includeErr != nil {
		fmt.Fprintf(w, "include: %s\n", cliColor(w, "error: "+includeErr.Error(), ansiRed))
	} else {
		fmt.Fprintf(w, "include: %s sha256=%s\n", filepath.Join(includeDir, "pawntest.inc"), includeHash)
	}

	if pawnccErr != nil {
		fmt.Fprintf(w, "pawncc: %s\n", cliColor(w, "not found: "+pawnccErr.Error(), ansiRed))
		fmt.Fprintf(w, "sample: %s\n", cliColor(w, "skipped", ansiDim))

		return errTestsFailed
	}

	fmt.Fprintf(w, "pawncc: %s (%s)\n", pawncc.Path, source)
	fmt.Fprintf(w, "pawncc_version: %s\n", version)

	if includeErr != nil {
		fmt.Fprintf(w, "sample: %s\n", cliColor(w, "skipped", ansiDim))
		return errTestsFailed
	}

	result := d.runDoctorSample(pawncc, includeDir)
	if result == "ok" {
		fmt.Fprintf(w, "sample: %s\n", cliColor(w, result, ansiGreen))
	} else {
		fmt.Fprintf(w, "sample: %s\n", cliColor(w, result, ansiRed))
	}

	if !versionOK || result != "ok" {
		return errTestsFailed
	}

	return nil
}

func (d DoctorCmd) resolveDoctorPawnCC(cacheDir string, explicit bool) (*compiler.Compiler, string, error) {
	if d.PawnCC != "" {
		path, err := resolveExecutable(d.PawnCC)
		if err != nil {
			return nil, "", err
		}

		if explicit {
			return compiler.FromPath(path), "explicit", nil
		}

		return compiler.FromPath(path), "config", nil
	}

	if path, err := exec.LookPath("pawncc"); err == nil {
		return compiler.FromPath(path), "PATH", nil
	}

	if c, ok := compiler.FindCachedCompiler(cacheDir); ok {
		return c, "cache", nil
	}

	return nil, "", compiler.ErrPawnCCNotFound
}

func resolveExecutable(path string) (string, error) {
	if strings.ContainsAny(path, `/\`) {
		if _, err := os.Stat(path); err != nil {
			return "", err
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			return path, nil
		}

		return abs, nil
	}

	return exec.LookPath(path)
}

func compilerVersion(c *compiler.Compiler) (string, bool) {
	cmd := c.Command()

	var out bytes.Buffer

	cmd.Stdout = &out

	cmd.Stderr = &out
	runErr := cmd.Run()

	for line := range strings.SplitSeq(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 3 && strings.EqualFold(fields[0], "Pawn") && strings.EqualFold(fields[1], "compiler") {
				return strings.Join(fields[:3], " "), true
			}

			return line, runErr == nil
		}
	}

	if runErr != nil {
		return "error: " + runErr.Error(), false
	}

	return "unknown", false
}

func (d DoctorCmd) runDoctorSample(c *compiler.Compiler, includeDir string) string {
	dir, err := os.MkdirTemp("", "pawntest-doctor-*")
	if err != nil {
		return "error: " + err.Error()
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "doctor.test.pwn")
	if err := os.WriteFile(src, []byte(`#include <pawntest>

TEST(doctor_sample)
{
    ASSERT_EQ(2 + 2, 4);
}
`), 0o644); err != nil {
		return "error: " + err.Error()
	}

	amx, err := compiler.Compile(src, compiler.Options{
		Compiler:  c,
		Includes:  append([]string{includeDir}, d.Include...),
		Defines:   d.Define,
		ExtraArgs: d.CompilerArg,
		OutDir:    dir,
		NoCache:   true,
		Count:     1,
	})
	if err != nil {
		return "compile error: " + err.Error()
	}

	suite, err := (runner.Runner{Backend: backend.NewGoAMXBackend()}).RunFile(amx)
	if err != nil {
		return "run error: " + err.Error()
	}

	if len(suite.Results) != 1 || suite.Results[0].Status != runner.Pass {
		return "failed"
	}

	return "ok"
}
