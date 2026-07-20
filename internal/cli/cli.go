package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"golang.org/x/sync/errgroup"

	pluginbridge "github.com/pawnkit/pawn-plugin-host/bridge"
	plugincontroller "github.com/pawnkit/pawn-plugin-host/controller"
	"github.com/pawnkit/pawntest/internal/backend"
	"github.com/pawnkit/pawntest/internal/cache"
	"github.com/pawnkit/pawntest/internal/compiler"
	"github.com/pawnkit/pawntest/internal/discovery"
	"github.com/pawnkit/pawntest/internal/report"
	"github.com/pawnkit/pawntest/internal/runner"
	watcher "github.com/pawnkit/pawntest/internal/watch"
)

var Version = "dev"

const (
	ExitOK       = 0
	ExitFailed   = 1
	ExitUsage    = 2
	ExitInternal = 3

	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
)

func cliColor(w io.Writer, s, code string) string {
	if !isTerminal(w) {
		return s
	}

	return code + s + ansiReset
}

func isTerminal(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	if force := os.Getenv("CLICOLOR_FORCE"); force != "" && force != "0" {
		return true
	}

	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	info, err := f.Stat()

	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

type SharedFlags struct {
	PawnCC      string   `name:"pawncc"      help:"Path to pawncc."`
	Include     []string `short:"i"           help:"Additional Pawn include directories."`
	Define      []string `short:"D"           help:"Compiler define passed to pawncc."`
	CompilerArg []string `name:"compiler-arg" help:"Extra raw pawncc argument."`
	CacheDir    string   `help:"Cache directory for compiled AMX files."`
}

type DoctorCmd struct {
	SharedFlags
	Provider []string `help:"Pawn native provider source or AMX file."`
}

type CacheCleanCmd struct {
	CacheDir string `help:"Cache directory to remove."`
}

type CacheCmd struct {
	Clean CacheCleanCmd `cmd:"" help:"Remove cached files."`
}

type TestCmd struct {
	Paths []string `arg:"" optional:"" name:"path" help:"Pawn test files, AMX files, or directories."`
	SharedFlags

	Run                 string        `help:"Only run/list tests whose name matches this regular expression."`
	Tags                string        `help:"Only run/list tests matching a tag expression, for example 'unit & !slow'."`
	List                bool          `help:"List tests instead of running them."`
	Recursive           bool          `help:"Recursively search input directories."`
	Format              Format        `default:"plain" enum:"plain,json,tap,junit" help:"Output format."`
	Color               Color         `default:"auto"  enum:"auto,always,never"   help:"Color plain output."`
	Output              string        `short:"o" help:"Write report to this file."`
	Verbose             bool          `short:"v" help:"Enable verbose output."`
	Quiet               bool          `short:"q" help:"Show only failures and the final summary."`
	NoCache             bool          `help:"Disable compile cache."`
	Count               int           `default:"0" help:"If set to 1, force recompilation similar to go test -count=1."`
	FailFast            bool          `help:"Stop after the first failing test."`
	AllowEmpty          bool          `help:"Exit successfully when no tests are found."`
	AllowUnknownNatives bool          `help:"Allow unconfigured unknown natives to return zero."`
	Isolation           string        `default:"test" enum:"test,suite" help:"Global-memory isolation between tests."`
	Shuffle             bool          `help:"Run tests in a reproducibly shuffled order."`
	Seed                int64         `default:"1" help:"Seed used by --shuffle."`
	Repeat              int           `default:"1" help:"Run each selected test this many times."`
	MaxInstructions     int           `default:"1000000" help:"Maximum AMX instructions per setup, test, or teardown invocation."`
	Jobs                int           `short:"j" default:"1" help:"Number of test files to compile and run concurrently."`
	UpdateSnapshots     bool          `help:"Create or replace golden string snapshots."`
	Coverage            bool          `help:"Collect source-line coverage."`
	CoverageOutput      string        `help:"Coverage output path (default: coverage.lcov or coverage.json)."`
	CoverageFormat      string        `default:"lcov" enum:"lcov,json" help:"Coverage output format."`
	Profile             bool          `help:"Collect instruction profile data."`
	ProfileOutput       string        `default:"profile.json" help:"Instruction profile output path."`
	Watch               bool          `help:"Rerun tests when Pawn sources, includes, or configuration change."`
	WatchInterval       time.Duration `default:"500ms" help:"Polling interval used by --watch."`
	FuzzSeed            int64         `default:"1" help:"Base seed for deterministic property tests."`
	Provider            []string      `help:"Pawn native provider source or AMX file."`
	NativePlugin        string        `help:"Legacy native plugin to run in an isolated worker."`
	PluginArchitecture  string        `default:"x86" enum:"x86,x64" help:"Native plugin architecture."`
	PluginWorker32      string        `help:"Path to the x86 plugin worker."`
	PluginWorker64      string        `help:"Path to the x64 plugin worker."`

	compilerCache *compiler.Compiler
	diagnostics   io.Writer

	stdinSrc    *os.File
	canPromptFn func() bool
	confirmFn   func() bool
	installFn   func(context.Context, string) (*compiler.Compiler, error)
}

type versionFlag bool

func (versionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)

	return nil
}

type CLI struct {
	Version versionFlag `name:"version" short:"V" help:"Print version and exit."`
	Test    TestCmd     `cmd:"" default:"withargs" help:"Run Pawn tests (default command)."`
	Doctor  DoctorCmd   `cmd:"doctor"              help:"Print environment diagnostics and run a sample compile/check."`
	Cache   CacheCmd    `cmd:"cache"               help:"Manage cached files."`
}

var errTestsFailed = errors.New("tests failed")

func Run(args []string, stdout, stderr io.Writer) int {
	var cli CLI

	parser, err := kong.New(
		&cli,
		kong.Name("pawntest"),
		kong.Description("Pawn test runner for SA-MP/open.mp-style projects."),
		kong.Vars{"version": Version},
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitInternal
	}

	parsed, err := parser.Parse(args)
	if err != nil {
		fmt.Fprintln(stderr, cliColor(stderr, err.Error(), ansiRed))
		return ExitUsage
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var execErr error

	switch parsed.Command() {
	case "doctor":
		execErr = cli.Doctor.execute(ctx, stdout)
	case "cache clean":
		execErr = cli.Cache.Clean.execute(stdout)
	default:
		execErr = cli.Test.execute(ctx, stdout, stderr)
	}

	if execErr != nil {
		if errors.Is(execErr, errTestsFailed) {
			return ExitFailed
		}

		fmt.Fprintln(stderr, cliColor(stderr, "error: "+execErr.Error(), ansiRed))

		return ExitUsage
	}

	return ExitOK
}

func (c CacheCleanCmd) execute(w io.Writer) error {
	dir := c.CacheDir
	if dir == "" {
		cfg, err := LoadDefaultConfig()
		if err != nil {
			return err
		}

		dir = cfg.CacheDir
		if dir == "" {
			dir = cache.Dir()
		}
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if abs == filepath.VolumeName(abs)+string(filepath.Separator) {
		return fmt.Errorf("refusing to remove filesystem root %s", abs)
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "removed %s\n", dir)

	return err
}

func (a TestCmd) compilerOption() *compiler.Compiler {
	if a.compilerCache != nil {
		return a.compilerCache
	}

	if a.PawnCC != "" {
		return compiler.FromPath(a.PawnCC)
	}

	return nil
}

func (a TestCmd) execute(ctx context.Context, stdout, stderr io.Writer) error {
	if a.Watch {
		return a.watch(ctx, stdout, stderr)
	}

	cfg, _, err := LoadDefaultConfigWithPath()
	if err != nil {
		return err
	}

	projectCfg, _, err := loadProjectConfig(workingDirectory())
	if err != nil {
		return err
	}

	projectCfg.merge(cfg)

	a.applyConfig(projectCfg)
	a.diagnostics = stderr

	if err := a.validate(); err != nil {
		return err
	}

	files, err := a.discoverFiles()
	if err != nil {
		return err
	}

	compilerInputs := append(append([]string{}, files...), a.Provider...)
	if err := a.ensureCompilerAvailable(ctx, compilerInputs, stderr); err != nil {
		return err
	}

	providers, err := a.ensureProviders(ctx)
	if err != nil {
		return err
	}

	coverage := a.newCoverage()
	profile := a.newProfile()

	r := a.newRunner(coverage, profile)

	pluginNatives, err := a.pluginNatives(ctx)
	if err != nil {
		return err
	}

	r.Natives = pluginNatives

	r.Providers = providers
	if a.List {
		return a.listTests(ctx, stdout, files, r)
	}

	runStarted := time.Now()

	suites, err := a.runFiles(ctx, files, r)
	if err != nil {
		return err
	}

	all := aggregateSuites(suites, time.Since(runStarted))
	if len(all.Results) == 0 && !a.AllowEmpty {
		return errors.New("no tests found")
	}

	if err := a.writeReportOutput(stdout, all); err != nil {
		return err
	}

	if coverage != nil {
		if err := a.writeCoverage(coverage); err != nil {
			return err
		}
	}

	if profile != nil {
		file, err := os.Create(a.ProfileOutput)
		if err != nil {
			return err
		}

		err = profile.WriteJSON(file)
		closeErr := file.Close()

		if err != nil {
			return err
		}

		if closeErr != nil {
			return closeErr
		}
	}

	if all.Failed() {
		return errTestsFailed
	}

	return nil
}

func (a TestCmd) validate() error {
	if a.Verbose && a.Quiet {
		return errors.New("verbose and quiet output cannot be enabled together")
	}

	return nil
}

func (a TestCmd) discoverFiles() ([]string, error) {
	files, err := discovery.Files(a.discoveryPaths())
	if err != nil {
		return nil, err
	}

	if len(files) == 0 && !a.AllowEmpty {
		return nil, errors.New("no test files found")
	}

	return files, nil
}

func (a TestCmd) newCoverage() *runner.Coverage {
	if !a.Coverage {
		return nil
	}

	return runner.NewCoverage()
}

func (a TestCmd) newProfile() *runner.Profile {
	if !a.Profile {
		return nil
	}

	return runner.NewProfile()
}

func (a TestCmd) pluginNatives(ctx context.Context) (map[string]backend.NativeFunc, error) {
	if a.NativePlugin == "" {
		return nil, nil
	}

	client := plugincontroller.Client{Worker32: a.PluginWorker32, Worker64: a.PluginWorker64}

	name, native, err := pluginbridge.Native(ctx, client, a.NativePlugin, a.PluginArchitecture)
	if err != nil {
		return nil, fmt.Errorf("native plugin: %w", err)
	}

	return map[string]backend.NativeFunc{name: native}, nil
}

func (a TestCmd) newRunner(coverage *runner.Coverage, profile *runner.Profile) runner.Runner {
	return runner.Runner{
		Backend:             backend.NewGoAMXBackend(),
		Run:                 a.Run,
		TagExpression:       a.Tags,
		FailFast:            a.FailFast,
		AllowUnknownNatives: a.AllowUnknownNatives,
		Isolation:           a.Isolation,
		Shuffle:             a.Shuffle,
		Seed:                a.Seed,
		Repeat:              a.Repeat,
		MaxInstructions:     a.MaxInstructions,
		Coverage:            coverage,
		Profile:             profile,
		FuzzSeed:            a.FuzzSeed,
	}
}

func (a TestCmd) ensureProviders(ctx context.Context) ([]string, error) {
	providers := make([]string, 0, len(a.Provider))
	for _, path := range a.Provider {
		provider, err := a.ensureAMX(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("compile provider %s: %w", path, err)
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

func (a TestCmd) listTests(ctx context.Context, stdout io.Writer, files []string, r runner.Runner) error {
	var names []string

	for _, file := range files {
		fileNames, err := a.listFileTests(ctx, file, r)
		if err != nil {
			return err
		}

		names = append(names, fileNames...)
	}

	if len(names) == 0 && !a.AllowEmpty {
		return errors.New("no tests found")
	}

	if a.Format == FormatJSON {
		return report.ListJSON(stdout, names)
	}

	return report.List(stdout, names)
}

func (a TestCmd) listFileTests(ctx context.Context, file string, r runner.Runner) ([]string, error) {
	expectations, err := compiler.DiagnosticExpectations(file)
	if err != nil {
		return nil, err
	}

	if len(expectations) > 0 {
		selected, selectErr := a.diagnosticSelected(file)
		if selectErr != nil {
			return nil, selectErr
		}

		if !selected {
			return nil, nil
		}

		return []string{diagnosticTestName(file)}, nil
	}

	amx, err := a.ensureAMX(ctx, file)
	if err != nil {
		return nil, err
	}

	tests, err := r.List(amx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(tests))
	for _, test := range tests {
		names = append(names, test.Name)
	}

	return names, nil
}

func aggregateSuites(suites []runner.Suite, duration time.Duration) runner.Suite {
	all := runner.Suite{Duration: duration}
	for _, suite := range suites {
		all.Results = append(all.Results, suite.Results...)
	}

	return all
}

func (a TestCmd) writeReportOutput(stdout io.Writer, suite runner.Suite) error {
	if a.Output == "" {
		return a.writeReport(stdout, suite)
	}

	f, err := os.Create(a.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	return a.writeReport(f, suite)
}

func (a TestCmd) watch(ctx context.Context, stdout, stderr io.Writer) error {
	cfg, _, err := LoadDefaultConfigWithPath()
	if err != nil {
		return err
	}

	a.applyConfig(cfg)

	a.Watch = false
	if len(a.Paths) == 0 {
		a.Paths = []string{"."}
	}

	paths := watcher.Paths(append(a.discoveryPaths(), a.Provider...), a.Include)
	snapshot := watcher.Files(paths)

	for {
		if err := a.execute(ctx, stdout, stderr); err != nil {
			fmt.Fprintln(stderr, err)
		}

		fmt.Fprintln(stdout, cliColor(stdout, "watching for changes...", ansiDim))

		next, err := watcher.Wait(ctx, paths, snapshot, a.WatchInterval)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		}

		snapshot = next

		fmt.Fprintln(stdout, cliColor(stdout, "\nchange detected; rerunning tests", ansiYellow))
	}
}

func (a TestCmd) writeCoverage(coverage *runner.Coverage) error {
	path := a.CoverageOutput
	if path == "" {
		path = "coverage." + a.CoverageFormat
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if a.CoverageFormat == "json" {
		return coverage.WriteJSON(file)
	}

	return coverage.WriteLCOV(file)
}

func (a TestCmd) runFiles(ctx context.Context, files []string, r runner.Runner) ([]runner.Suite, error) {
	jobs := a.Jobs
	if jobs < 1 {
		return nil, errors.New("jobs must be at least 1")
	}

	if a.FailFast {
		suites := make([]runner.Suite, 0, len(files))
		for _, file := range files {
			suite, err := a.runFile(ctx, file, r)
			if err != nil {
				return nil, err
			}

			suites = append(suites, suite)
			if suite.Failed() {
				break
			}
		}

		return suites, nil
	}

	type fileResult struct {
		suite runner.Suite
		err   error
	}

	results := make([]fileResult, len(files))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(jobs)

	for i, file := range files {
		g.Go(func() error {
			results[i].suite, results[i].err = a.runFile(ctx, file, r)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	suites := make([]runner.Suite, 0, len(results))
	for _, result := range results {
		if result.err != nil {
			return nil, result.err
		}

		suites = append(suites, result.suite)
	}

	return suites, nil
}

func (a TestCmd) runFile(ctx context.Context, file string, r runner.Runner) (runner.Suite, error) {
	expectations, err := compiler.DiagnosticExpectations(file)

	var suite runner.Suite

	if err == nil && len(expectations) > 0 {
		var selected bool

		selected, err = a.diagnosticSelected(file)
		if err == nil && !selected {
			err = runner.ErrNoTestsFound
		} else if err == nil {
			suite, err = a.runDiagnosticTest(ctx, file, expectations)
		}
	} else if err == nil {
		var amx string

		amx, err = a.ensureAMX(ctx, file)
		if err == nil {
			r.SourcePath = file
			r.UpdateSnapshots = a.UpdateSnapshots
			suite, err = r.RunFile(amx)
		}
	}

	if a.allowsEmptyFileSelection() && errors.Is(err, runner.ErrNoTestsFound) {
		err = nil
	}

	for i := range suite.Results {
		suite.Results[i].Source = file
	}

	return suite, err
}

func (a TestCmd) allowsEmptyFileSelection() bool {
	return a.AllowEmpty || a.Run != "" || a.Tags != ""
}

func (a TestCmd) runDiagnosticTest(ctx context.Context, path string, expectations []compiler.DiagnosticExpectation) (runner.Suite, error) {
	opts, err := a.compilerOptions(ctx, path)
	if err != nil {
		return runner.Suite{}, err
	}

	result, err := compiler.CheckDiagnostics(path, opts, expectations)
	if err != nil {
		return runner.Suite{}, err
	}

	status := runner.Fail
	if result.Passed {
		status = runner.Pass
	}

	return runner.Suite{Results: []runner.Result{{
		Name: diagnosticTestName(path), File: path, Source: path, Status: status, Message: result.Message,
	}}}, nil
}

func diagnosticTestName(path string) string {
	return "compile:" + strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func (a TestCmd) diagnosticSelected(path string) (bool, error) {
	matchesTags, err := runner.MatchTags(a.Tags)
	if err != nil || !matchesTags {
		return false, err
	}

	if a.Run == "" {
		return true, nil
	}

	expression, err := regexp.Compile(a.Run)
	if err != nil {
		return false, fmt.Errorf("invalid run regex: %w", err)
	}

	return expression.MatchString(diagnosticTestName(path)), nil
}

func (a *TestCmd) ensureCompilerAvailable(ctx context.Context, files []string, stderr io.Writer) error {
	if a.PawnCC != "" || !needsCompiler(files) {
		return nil
	}

	if _, err := exec.LookPath("pawncc"); err == nil {
		return nil
	}

	cacheDir := a.CacheDir
	if cacheDir == "" {
		cacheDir = cache.Dir()
	}

	if c, ok := compiler.FindCachedCompiler(cacheDir); ok {
		a.PawnCC = c.Path
		a.compilerCache = c

		return nil
	}

	if !a.resolveCanPrompt() {
		return compiler.ErrPawnCCNotFound
	}

	fmt.Fprintf(stderr, "pawncc was not found in PATH. Download the openmultiplayer compiler from GitHub releases to %s? [y/N] ", filepath.Join(cacheDir, "tools", "openmp-compiler"))

	if !a.resolveConfirm() {
		return compiler.ErrPawnCCNotFound
	}

	fmt.Fprintln(stderr, cliColor(stderr, "Downloading openmultiplayer compiler...", ansiCyan))

	c, err := a.resolveInstall(ctx, cacheDir)
	if err != nil {
		return err
	}

	a.PawnCC = c.Path
	a.compilerCache = c
	fmt.Fprintln(stderr, cliColor(stderr, "Using downloaded compiler: "+c.Path, ansiGreen))

	return nil
}

func needsCompiler(files []string) bool {
	for _, file := range files {
		if !strings.EqualFold(filepath.Ext(file), ".amx") {
			return true
		}
	}

	return false
}

func (a *TestCmd) resolveStdin() *os.File {
	if a.stdinSrc != nil {
		return a.stdinSrc
	}

	return os.Stdin
}

func (a *TestCmd) resolveCanPrompt() bool {
	if a.canPromptFn != nil {
		return a.canPromptFn()
	}

	return canPrompt(a.resolveStdin())
}

func (a *TestCmd) resolveConfirm() bool {
	if a.confirmFn != nil {
		return a.confirmFn()
	}

	return confirmDownload(a.resolveStdin())
}

func (a *TestCmd) resolveInstall(ctx context.Context, dir string) (*compiler.Compiler, error) {
	if a.installFn != nil {
		return a.installFn(ctx, dir)
	}

	return compiler.InstallOpenMPCompiler(ctx, dir)
}

func canPrompt(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func confirmDownload(r io.Reader) bool {
	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

func (a TestCmd) ensureAMX(ctx context.Context, path string) (string, error) {
	if filepath.Ext(path) == ".amx" {
		return path, nil
	}

	opts, err := a.compilerOptions(ctx, path)
	if err != nil {
		return "", err
	}

	return compiler.CompileContext(ctx, path, opts)
}

func (a TestCmd) compilerOptions(ctx context.Context, path string) (compiler.Options, error) {
	_ = ctx

	cacheDir := a.CacheDir
	if cacheDir == "" {
		cacheDir = cache.Dir()
	}

	includeDir, err := cache.IncludeDirIn(cacheDir)
	if err != nil {
		return compiler.Options{}, err
	}

	includes := append([]string{includeDir}, a.Include...)

	outDir, err := cache.AMXDirIn(cacheDir)
	if err != nil {
		return compiler.Options{}, err
	}

	return compiler.Options{
		Compiler:    a.compilerOption(),
		Includes:    includes,
		Defines:     a.Define,
		ExtraArgs:   a.CompilerArg,
		OutDir:      outDir,
		NoCache:     a.NoCache,
		Count:       a.Count,
		Diagnostics: a.diagnostics,
	}, nil
}

func (a TestCmd) discoveryPaths() []string {
	if !a.Recursive {
		return a.Paths
	}

	paths := a.Paths
	if len(paths) == 0 {
		paths = []string{"."}
	}

	out := make([]string, 0, len(paths))
	for _, path := range paths {
		if filepath.Ext(path) == "" && !filepath.IsAbs(path) {
			out = append(out, path+"/...")
			continue
		}

		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			out = append(out, path+"/...")
			continue
		}

		out = append(out, path)
	}

	return out
}

func (a TestCmd) writeReport(w io.Writer, suite runner.Suite) error {
	switch a.Format {
	case FormatPlain:
		return report.PlainWithOptions(w, suite, report.PlainOptions{Color: a.colorEnabled(w), Verbose: a.Verbose, Quiet: a.Quiet})
	case FormatJSON:
		return report.JSON(w, suite)
	case FormatTAP:
		return report.TAP(w, suite)
	case FormatJUnit:
		return report.JUnit(w, suite)
	default:
		return fmt.Errorf("invalid format %q", a.Format)
	}
}

func (a TestCmd) colorEnabled(w io.Writer) bool {
	switch a.Color {
	case ColorAuto:
		return isTerminal(w)
	case ColorAlways:
		return true
	case ColorNever:
		return false
	}

	return false
}
