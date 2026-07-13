package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type (
	Format string
	Color  string
)

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
	FormatTAP   Format = "tap"
	FormatJUnit Format = "junit"

	ColorAuto   Color = "auto"
	ColorAlways Color = "always"
	ColorNever  Color = "never"
)

type Config struct {
	PawnCC              string   `json:"pawncc"               yaml:"pawncc"               toml:"pawncc"`
	Include             []string `json:"include"              yaml:"include"               toml:"include"`
	Tests               []string `json:"tests"                yaml:"tests"                toml:"tests"`
	Format              Format   `json:"format"               yaml:"format"               toml:"format"`
	CacheDir            string   `json:"cache_dir"            yaml:"cache_dir"            toml:"cache_dir"`
	AllowUnknownNatives bool     `json:"allow_unknown_natives" yaml:"allow_unknown_natives" toml:"allow_unknown_natives"`
	Isolation           string   `json:"isolation"            yaml:"isolation"            toml:"isolation"`
	Run                 string   `json:"run"                  yaml:"run"                  toml:"run"`
	Tags                string   `json:"tags"                 yaml:"tags"                 toml:"tags"`
	Shuffle             bool     `json:"shuffle"              yaml:"shuffle"              toml:"shuffle"`
	Seed                int64    `json:"seed"                 yaml:"seed"                 toml:"seed"`
	Repeat              int      `json:"repeat"               yaml:"repeat"               toml:"repeat"`
	MaxInstructions     int      `json:"max_instructions"     yaml:"max_instructions"     toml:"max_instructions"`
	Jobs                int      `json:"jobs"                 yaml:"jobs"                 toml:"jobs"`
	UpdateSnapshots     bool     `json:"update_snapshots"     yaml:"update_snapshots"     toml:"update_snapshots"`
	Coverage            bool     `json:"coverage"             yaml:"coverage"             toml:"coverage"`
	CoverageOutput      string   `json:"coverage_output"      yaml:"coverage_output"      toml:"coverage_output"`
	CoverageFormat      string   `json:"coverage_format"      yaml:"coverage_format"      toml:"coverage_format"`
	FuzzSeed            int64    `json:"fuzz_seed"            yaml:"fuzz_seed"            toml:"fuzz_seed"`
	AllowEmpty          bool     `json:"allow_empty"          yaml:"allow_empty"          toml:"allow_empty"`
	Verbose             bool     `json:"verbose"              yaml:"verbose"              toml:"verbose"`
	Quiet               bool     `json:"quiet"                yaml:"quiet"                toml:"quiet"`
	Providers           []string `json:"providers"            yaml:"providers"            toml:"providers"`
}

func LoadDefaultConfig() (Config, error) {
	cfg, _, err := LoadDefaultConfigWithPath()
	return cfg, err
}

func LoadDefaultConfigWithPath() (Config, string, error) {
	for _, path := range []string{"pawntest.json", "pawntest.yaml", "pawntest.yml", "pawntest.toml"} {
		cfg, found, err := loadConfigIfExists(path)
		if err != nil || found {
			return cfg, path, err
		}
	}

	return Config{}, "", nil
}

func LoadConfig(path string) (Config, error) {
	cfg, _, err := loadConfigIfExists(path)
	return cfg, err
}

func loadConfigIfExists(path string) (Config, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, false, nil
		}

		return Config{}, false, err
	}

	var cfg Config

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		decoder := yaml.NewDecoder(bytes.NewReader(data))
		decoder.KnownFields(true)

		if err := decoder.Decode(&cfg); err != nil {
			return Config{}, true, err
		}
	case ".toml":
		metadata, err := toml.Decode(string(data), &cfg)
		if err != nil {
			return Config{}, true, err
		}

		if undecoded := metadata.Undecoded(); len(undecoded) > 0 {
			return Config{}, true, fmt.Errorf("unknown config field %q", undecoded[0])
		}
	default:
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&cfg); err != nil {
			return Config{}, true, err
		}
	}

	if err := cfg.validate(); err != nil {
		return Config{}, true, err
	}

	return cfg, true, nil
}

func (cfg Config) validate() error {
	switch cfg.Format {
	case "", FormatPlain, FormatJSON, FormatTAP, FormatJUnit:
	default:
		return fmt.Errorf("invalid config format %q", cfg.Format)
	}

	if cfg.Isolation != "" && cfg.Isolation != "test" && cfg.Isolation != "suite" {
		return fmt.Errorf("invalid isolation %q", cfg.Isolation)
	}

	if cfg.CoverageFormat != "" && cfg.CoverageFormat != "lcov" && cfg.CoverageFormat != "json" {
		return fmt.Errorf("invalid coverage format %q", cfg.CoverageFormat)
	}

	if cfg.Jobs < 0 || cfg.Repeat < 0 || cfg.MaxInstructions < 0 {
		return errors.New("jobs, repeat, and max_instructions cannot be negative")
	}

	return nil
}

func (a *TestCmd) applyConfig(cfg Config) {
	if len(a.Paths) == 0 {
		a.Paths = cfg.Tests
	}

	if a.PawnCC == "" {
		a.PawnCC = cfg.PawnCC
	}

	if len(a.Include) == 0 {
		a.Include = cfg.Include
	}

	if a.Format == "" || a.Format == FormatPlain {
		if cfg.Format != "" {
			a.Format = cfg.Format
		}
	}

	if a.CacheDir == "" {
		a.CacheDir = cfg.CacheDir
	}

	if cfg.AllowUnknownNatives {
		a.AllowUnknownNatives = true
	}

	if a.Isolation == "test" && cfg.Isolation != "" {
		a.Isolation = cfg.Isolation
	}

	if cfg.AllowEmpty {
		a.AllowEmpty = true
	}

	if cfg.Verbose {
		a.Verbose = true
	}

	if cfg.Quiet {
		a.Quiet = true
	}

	if a.Run == "" {
		a.Run = cfg.Run
	}

	if a.Tags == "" {
		a.Tags = cfg.Tags
	}

	if cfg.Shuffle {
		a.Shuffle = true
	}

	if a.Seed == 1 && cfg.Seed != 0 {
		a.Seed = cfg.Seed
	}

	if a.Repeat == 1 && cfg.Repeat > 0 {
		a.Repeat = cfg.Repeat
	}

	if a.MaxInstructions == 1_000_000 && cfg.MaxInstructions > 0 {
		a.MaxInstructions = cfg.MaxInstructions
	}

	if a.Jobs == 1 && cfg.Jobs > 0 {
		a.Jobs = cfg.Jobs
	}

	if cfg.UpdateSnapshots {
		a.UpdateSnapshots = true
	}

	if cfg.Coverage {
		a.Coverage = true
	}

	if a.CoverageOutput == "" {
		a.CoverageOutput = cfg.CoverageOutput
	}

	if a.CoverageFormat == "lcov" && cfg.CoverageFormat != "" {
		a.CoverageFormat = cfg.CoverageFormat
	}

	if a.FuzzSeed == 1 && cfg.FuzzSeed != 0 {
		a.FuzzSeed = cfg.FuzzSeed
	}

	if len(a.Provider) == 0 {
		a.Provider = cfg.Providers
	}
}

func (d *DoctorCmd) applyConfig(cfg Config) {
	if d.PawnCC == "" {
		d.PawnCC = cfg.PawnCC
	}

	if len(d.Include) == 0 {
		d.Include = cfg.Include
	}

	if d.CacheDir == "" {
		d.CacheDir = cfg.CacheDir
	}

	if len(d.Provider) == 0 {
		d.Provider = cfg.Providers
	}
}
