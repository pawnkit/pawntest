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
	_, err := testLocations(path, true)
	return err
}

func TestLocations(path string) (map[string]int, error) {
	return testLocations(path, false)
}

func testLocations(path string, validate bool) (map[string]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	locations := make(map[string]int)
	scanner := bufio.NewScanner(f)
	line := 0

	for scanner.Scan() {
		line++

		match := generatedSymbolPattern.FindStringSubmatch(stripLineComment(scanner.Text()))
		if match == nil {
			continue
		}

		macro, name := match[1], match[2]
		if macro != "FIXTURE" {
			locations["test_"+name] = line
		}

		prefix, suffix := "test_", ""
		if macro == "FIXTURE" {
			prefix, suffix = "fixture_", "_teardown"
		}

		generated := prefix + name + suffix
		if !validate || len(generated) <= pawnSymbolLimit {
			continue
		}

		maxName := pawnSymbolLimit - len(prefix) - len(suffix)

		return nil, fmt.Errorf(
			"%s:%d: %s name %q generates Pawn symbol %q (%d characters); maximum source name length is %d",
			path, line, macro, name, generated, len(generated), maxName,
		)
	}

	return locations, scanner.Err()
}

func stripLineComment(line string) string {
	if before, _, ok := strings.Cut(line, "//"); ok {
		return before
	}

	return line
}
