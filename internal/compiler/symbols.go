package compiler

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const pawnSymbolLimit = 31

var generatedSymbolPattern = regexp.MustCompile(`^\s*(TEST(?:_CASE[23]?)?|FIXTURE)\s*\(\s*([A-Za-z_][A-Za-z0-9_]*)`)

func validateGeneratedSymbols(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0

	for scanner.Scan() {
		line++

		match := generatedSymbolPattern.FindStringSubmatch(stripLineComment(scanner.Text()))
		if match == nil {
			continue
		}

		macro, name := match[1], match[2]

		prefix, suffix := "test_", ""
		if macro == "FIXTURE" {
			prefix, suffix = "fixture_", "_teardown"
		}

		generated := prefix + name + suffix
		if len(generated) <= pawnSymbolLimit {
			continue
		}

		maxName := pawnSymbolLimit - len(prefix) - len(suffix)

		return fmt.Errorf(
			"%s:%d: %s name %q generates Pawn symbol %q (%d characters); maximum source name length is %d",
			path, line, macro, name, generated, len(generated), maxName,
		)
	}

	return scanner.Err()
}

func stripLineComment(line string) string {
	if before, _, ok := strings.Cut(line, "//"); ok {
		return before
	}

	return line
}
