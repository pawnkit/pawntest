package report

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pawnkit/pawntest/internal/runner"
)

func Plain(w io.Writer, suite runner.Suite) error {
	return PlainWithOptions(w, suite, PlainOptions{})
}

type PlainOptions struct {
	Color   bool
	Verbose bool
	Quiet   bool
}

func PlainWithOptions(w io.Writer, suite runner.Suite, opts PlainOptions) error {
	groups := groupBySource(suite.Results)
	printedGroup := false
	for _, g := range groups {
		visible := visibleResults(g.results, opts.Quiet)
		if len(visible) == 0 {
			continue
		}
		if printedGroup {
			fmt.Fprintln(w)
		}
		printedGroup = true
		if g.source != "" {
			fmt.Fprintln(w, displaySource(g.source, opts.Verbose))
		}
		for _, r := range visible {
			writePlainResult(w, r, g.source, opts)
		}
	}
	if printedGroup {
		fmt.Fprintln(w)
	}
	writePlainSummary(w, suite, len(groups), opts.Color)
	return nil
}

func writePlainResult(w io.Writer, result runner.Result, source string, opts PlainOptions) {
	indent := resultIndent(source)
	duration := ""
	if opts.Verbose {
		duration = "  " + colorize(formatDuration(result.Duration), ansiDim, opts.Color)
	}

	fmt.Fprintf(w, "%s%s %s%s\n", indent, formatStatus(result.Status, opts.Color), result.Name, duration)
	writeIndentedLines(w, indent, result.Message)

	if location := formatLocation(result, source, opts.Verbose); location != "" {
		fmt.Fprintf(w, "%s      %s\n", indent, colorize(location, ansiDim, opts.Color))
	}
	if command := failureRerunCommand(result, source); command != "" {
		fmt.Fprintf(w, "%s      %s\n", indent, colorize("rerun: "+command, ansiDim, opts.Color))
	}
}

func resultIndent(source string) string {
	if source == "" {
		return ""
	}
	return "  "
}

func writeIndentedLines(w io.Writer, indent, text string) {
	if text == "" {
		return
	}
	for _, line := range strings.Split(text, "\n") {
		fmt.Fprintf(w, "%s      %s\n", indent, line)
	}
}

func failureRerunCommand(result runner.Result, source string) string {
	if !isFailedStatus(result.Status) {
		return ""
	}
	return rerunCommand(result, source)
}

func visibleResults(results []runner.Result, quiet bool) []runner.Result {
	if !quiet {
		return results
	}
	visible := make([]runner.Result, 0, len(results))
	for _, result := range results {
		if isFailedStatus(result.Status) {
			visible = append(visible, result)
		}
	}
	return visible
}

func writePlainSummary(w io.Writer, suite runner.Suite, files int, color bool) {
	summary := suite.Summary()
	overall, overallColor := "PASS", ansiGreen
	if suite.Failed() {
		overall, overallColor = "FAIL", ansiRed
	}
	testWord, fileWord := plural(summary.Total, "test", "tests"), plural(files, "file", "files")
	duration := suite.Duration
	if duration == 0 {
		for _, result := range suite.Results {
			duration += result.Duration
		}
	}
	fmt.Fprintf(w, "%s  %d %s across %d %s in %s\n", colorize(overall, overallColor, color), summary.Total, testWord, files, fileWord, formatDuration(duration))

	counts := []string{colorize(fmt.Sprintf("%d passed", summary.Passed), ansiGreen, color)}
	if summary.Failed > 0 {
		counts = append(counts, colorize(fmt.Sprintf("%d failed", summary.Failed), ansiRed, color))
	}
	if summary.Skipped > 0 {
		counts = append(counts, colorize(fmt.Sprintf("%d skipped", summary.Skipped), ansiYellow, color))
	}
	if summary.Errored > 0 {
		counts = append(counts, colorize(fmt.Sprintf("%d errored", summary.Errored), ansiRed, color))
	}
	if summary.XFailed > 0 {
		counts = append(counts, colorize(fmt.Sprintf("%d xfailed", summary.XFailed), ansiYellow, color))
	}
	if summary.XPassed > 0 {
		counts = append(counts, colorize(fmt.Sprintf("%d xpassed", summary.XPassed), ansiRed, color))
	}
	fmt.Fprintf(w, "      %s\n", strings.Join(counts, ", "))
}

func displaySource(source string, verbose bool) string {
	if !verbose {
		return source
	}
	absolute, err := filepath.Abs(source)
	if err == nil {
		return absolute
	}
	return source
}

func formatLocation(result runner.Result, source string, verbose bool) string {
	if result.Status == runner.Pass || result.File == "" {
		return ""
	}
	if !verbose && source != "" && samePath(result.File, source) {
		if result.Line > 0 {
			return fmt.Sprintf("at line %d", result.Line)
		}
		return ""
	}
	file := result.File
	if verbose {
		if absolute, err := filepath.Abs(file); err == nil {
			file = absolute
		}
	}
	if result.Line > 0 {
		return fmt.Sprintf("at %s:%d", file, result.Line)
	}
	return "at " + file
}

func samePath(left, right string) bool {
	leftPath, leftErr := filepath.Abs(left)
	rightPath, rightErr := filepath.Abs(right)
	return leftErr == nil && rightErr == nil && filepath.Clean(leftPath) == filepath.Clean(rightPath)
}

func rerunCommand(result runner.Result, source string) string {
	if source == "" {
		source = result.Source
	}
	if source == "" {
		return ""
	}
	name := result.Name
	if index := strings.Index(name, " [attempt "); index >= 0 {
		name = name[:index]
	}
	return "pawntest " + shellQuote(source) + " --run " + shellQuote("^"+regexp.QuoteMeta(name)+"$")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func formatDuration(duration time.Duration) string {
	if duration <= 0 {
		return "0ms"
	}
	if duration < time.Millisecond {
		return "<1ms"
	}
	return duration.Round(time.Millisecond).String()
}

func plural(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func isFailedStatus(status runner.Status) bool {
	return status == runner.Fail || status == runner.Error || status == runner.XPass
}

type fileGroup struct {
	source  string
	results []runner.Result
}

func groupBySource(results []runner.Result) []fileGroup {
	if len(results) == 0 {
		return nil
	}
	orders := make([]string, 0, len(results))
	groups := make(map[string][]runner.Result)
	for _, r := range results {
		src := r.Source
		if src == "" {
			src = "__nosource"
		}
		if _, ok := groups[src]; !ok {
			orders = append(orders, src)
		}
		groups[src] = append(groups[src], r)
	}
	out := make([]fileGroup, 0, len(orders))
	for _, src := range orders {
		if src == "__nosource" {
			out = append(out, fileGroup{results: groups[src]})
		} else {
			out = append(out, fileGroup{source: src, results: groups[src]})
		}
	}
	return out
}

func List(w io.Writer, tests []string) error {
	for _, t := range tests {
		fmt.Fprintln(w, t)
	}
	return nil
}

const (
	ansiReset  = "\x1b[0m"
	ansiDim    = "\x1b[2m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
)

func formatStatus(status runner.Status, color bool) string {
	label := fmt.Sprintf("%-5s", strings.ToUpper(string(status)))
	switch status {
	case runner.Pass:
		return colorize(label, ansiGreen, color)
	case runner.Fail, runner.Error:
		return colorize(label, ansiRed, color)
	case runner.Skip:
		return colorize(label, ansiYellow, color)
	case runner.XFail:
		return colorize(label, ansiYellow, color)
	case runner.XPass:
		return colorize(label, ansiRed, color)
	default:
		return label
	}
}

func colorize(s, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return code + s + ansiReset
}
