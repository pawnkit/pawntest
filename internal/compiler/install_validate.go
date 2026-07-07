package compiler

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

func validateCompiler(c *Compiler) error {
	cmd := c.Command()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		if compilerStarted(out.String()) {
			return nil
		}
		details := compilerValidationDetails(out.String(), runtime.GOOS)
		return fmt.Errorf("downloaded pawncc could not run: %w%s", err, details)
	}
	return nil
}

func compilerStarted(output string) bool {
	return strings.Contains(output, "Pawn compiler") && strings.Contains(output, "Usage:")
}

func compilerValidationDetails(output, goos string) string {
	var b strings.Builder

	output = strings.TrimSpace(output)
	if output != "" {
		b.WriteString("\n")
		b.WriteString(output)
	}
	if hint := compilerRuntimeHint(goos); hint != "" {
		b.WriteString("\n")
		b.WriteString(hint)
	}
	if b.Len() == 0 {
		return ""
	}
	return ": " + b.String()
}

func compilerRuntimeHint(goos string) string {
	if goos != "linux" {
		return ""
	}
	return "hint: the official openmultiplayer Linux compiler is a 32-bit glibc binary and needs the runtime loader/libraries visible where pawntest is running, including /lib/ld-linux.so.2. On Arch/CachyOS install lib32-glibc; on Debian/Ubuntu install libc6-i386."
}
