package compiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const defaultPawnCC = "pawncc"

var (
	diagnosticDirective = regexp.MustCompile(`(?im)^\s*//\s*pawntest:\s*expect-(error|warning)\s+(\d{3})\s*$`)
	compilerDiagnostic  = regexp.MustCompile(`(?i)\b(error|warning)\s+(\d{3}):`)
	pawntestInclude     = regexp.MustCompile(`(?im)^\s*#include\s*[<"]pawntest(?:\.inc)?[>"]`)
)

type DiagnosticExpectation struct {
	Kind string
	Code string
}

type DiagnosticResult struct {
	Passed  bool
	Message string
}

func DiagnosticExpectations(path string) ([]DiagnosticExpectation, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	matches := diagnosticDirective.FindAllStringSubmatch(string(source), -1)
	if len(matches) > 0 && !pawntestInclude.Match(source) {
		return nil, fmt.Errorf("%s must include <pawntest>", path)
	}

	out := make([]DiagnosticExpectation, 0, len(matches))
	for _, match := range matches {
		out = append(out, DiagnosticExpectation{Kind: strings.ToLower(match[1]), Code: match[2]})
	}

	return out, nil
}

func CheckDiagnostics(path string, opts Options, expected []DiagnosticExpectation) (DiagnosticResult, error) {
	if len(expected) == 0 {
		return DiagnosticResult{}, fmt.Errorf("no compiler diagnostic expectations in %s", path)
	}

	tempDir, err := os.MkdirTemp("", "pawntest-diagnostics-")
	if err != nil {
		return DiagnosticResult{}, err
	}
	defer os.RemoveAll(tempDir)

	c := opts.Compiler

	pawnccPath := defaultPawnCC
	if c != nil && c.Path != "" {
		pawnccPath = c.Path
	}

	args := buildArgs(filepath.Join(tempDir, "diagnostic.amx"), path, opts)

	var cmd *exec.Cmd
	if c != nil {
		cmd = c.Command(args...)
	} else {
		cmd = exec.Command(pawnccPath, args...)
	}

	var output bytes.Buffer

	cmd.Stdout, cmd.Stderr = &output, &output

	runErr := cmd.Run()
	if runErr != nil && pawnccPath == defaultPawnCC {
		if _, lookErr := exec.LookPath(pawnccPath); lookErr != nil {
			return DiagnosticResult{}, ErrPawnCCNotFound
		}
	}

	found := map[string]bool{}
	for _, match := range compilerDiagnostic.FindAllStringSubmatch(output.String(), -1) {
		found[strings.ToLower(match[1])+":"+match[2]] = true
	}

	var missing []string

	expectedSet := map[string]bool{}

	for _, expectation := range expected {
		key := expectation.Kind + ":" + expectation.Code

		expectedSet[key] = true
		if !found[key] {
			missing = append(missing, expectation.Kind+" "+expectation.Code)
		}
	}

	var unexpected []string

	for diagnostic := range found {
		if strings.HasPrefix(diagnostic, "error:") && !expectedSet[diagnostic] {
			unexpected = append(unexpected, strings.Replace(diagnostic, ":", " ", 1))
		}
	}

	if len(unexpected) > 0 {
		sort.Strings(unexpected)
		missing = append(missing, "unexpected "+strings.Join(unexpected, ", "))
	}

	if len(missing) > 0 {
		sort.Strings(missing)

		message := "compiler diagnostic mismatch: " + strings.Join(missing, ", ")
		if text := strings.TrimSpace(output.String()); text != "" {
			message += "\n" + text
		}

		return DiagnosticResult{Message: message}, nil
	}

	return DiagnosticResult{Passed: true}, nil
}
