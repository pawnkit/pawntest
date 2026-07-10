package pawntest

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
	"github.com/pawnkit/pawntest/internal/cache"
	"github.com/pawnkit/pawntest/internal/compiler"
	"github.com/pawnkit/pawntest/internal/runner"
)

type Runner struct {
	Run                 string
	FailFast            bool
	AllowUnknownNatives bool
	PawnCC              string
	Include             []string
	Define              []string
	CompilerArg         []string
	CacheDir            string
	NoCache             bool
	Count               int
	Isolation           string
	Tags                string
	Shuffle             bool
	Seed                int64
	Repeat              int
	MaxInstructions     int
	Natives             map[string]NativeFunc
	UpdateSnapshots     bool
	FuzzSeed            int64
}

func (r Runner) List(path string) ([]Public, error) {
	amx, err := r.ensureAMX(path)
	if err != nil {
		return nil, err
	}

	publics, err := r.internal().List(amx)
	if err != nil {
		return nil, err
	}

	out := make([]Public, 0, len(publics))
	for _, pub := range publics {
		out = append(out, Public{Index: pub.Index, Name: pub.Name})
	}

	return out, nil
}

func (r Runner) RunFile(path string) (Suite, error) {
	amx, err := r.ensureAMX(path)
	if err != nil {
		return Suite{}, err
	}

	internal := r.internal()
	internal.SourcePath = path
	internal.UpdateSnapshots = r.UpdateSnapshots

	suite, err := internal.RunFile(amx)
	if err != nil {
		return Suite{}, err
	}

	out := Suite{Results: make([]Result, 0, len(suite.Results))}
	for _, result := range suite.Results {
		out.Results = append(out.Results, Result{
			Name:     result.Name,
			Source:   result.Source,
			File:     result.File,
			Line:     result.Line,
			Status:   Status(result.Status),
			Message:  result.Message,
			Duration: result.Duration.Milliseconds(),
		})
	}

	return out, nil
}

func (r Runner) internal() runner.Runner {
	natives := make(map[string]backend.NativeFunc, len(r.Natives))
	for name, native := range r.Natives {
		fn := native
		natives[name] = func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
			publicParams := make([]Cell, len(params))
			for i, param := range params {
				publicParams[i] = Cell(param)
			}

			value, err := fn(nativeContext{NativeContext: ctx}, publicParams)

			return backend.Cell(value), err
		}
	}

	return runner.Runner{
		Run:                 r.Run,
		FailFast:            r.FailFast,
		AllowUnknownNatives: r.AllowUnknownNatives,
		Isolation:           r.Isolation,
		TagExpression:       r.Tags,
		Shuffle:             r.Shuffle,
		Seed:                r.Seed,
		Repeat:              r.Repeat,
		MaxInstructions:     r.MaxInstructions,
		Natives:             natives,
		FuzzSeed:            r.FuzzSeed,
	}
}

type Cell int32

type NativeContext interface {
	ReadString(addr Cell) (string, error)
	WriteString(addr Cell, value string) error
	ReadCell(addr Cell) (Cell, error)
	WriteCell(addr, value Cell) error
	CallPublic(name string, args ...Cell) (Cell, error)
}

type NativeFunc func(ctx NativeContext, params []Cell) (Cell, error)

type nativeContext struct{ backend.NativeContext }

func (ctx nativeContext) ReadString(addr Cell) (string, error) {
	return ctx.NativeContext.ReadString(backend.Cell(addr))
}

func (ctx nativeContext) WriteString(addr Cell, value string) error {
	return ctx.NativeContext.WriteString(backend.Cell(addr), value)
}

func (ctx nativeContext) ReadCell(addr Cell) (Cell, error) {
	value, err := ctx.NativeContext.ReadCell(backend.Cell(addr))
	return Cell(value), err
}

func (ctx nativeContext) WriteCell(addr, value Cell) error {
	return ctx.NativeContext.WriteCell(backend.Cell(addr), backend.Cell(value))
}

func (ctx nativeContext) CallPublic(name string, args ...Cell) (Cell, error) {
	caller, ok := ctx.NativeContext.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support public callbacks")
	}

	internalArgs := make([]backend.Cell, len(args))
	for i, arg := range args {
		internalArgs[i] = backend.Cell(arg)
	}

	value, err := caller.CallPublic(name, internalArgs...)

	return Cell(value), err
}

func (r Runner) ensureAMX(path string) (string, error) {
	if strings.EqualFold(filepath.Ext(path), ".amx") {
		return path, nil
	}

	cacheDir := r.CacheDir
	if cacheDir == "" {
		cacheDir = cache.Dir()
	}

	includeDir, err := cache.IncludeDirIn(cacheDir)
	if err != nil {
		return "", err
	}

	outDir, err := cache.AMXDirIn(cacheDir)
	if err != nil {
		return "", err
	}

	includes := append([]string{includeDir}, r.Include...)

	var comp *compiler.Compiler
	if r.PawnCC != "" {
		comp = compiler.Bare(r.PawnCC)
	}

	return compiler.Compile(path, compiler.Options{
		Compiler:  comp,
		Includes:  includes,
		Defines:   r.Define,
		ExtraArgs: r.CompilerArg,
		OutDir:    outDir,
		NoCache:   r.NoCache,
		Count:     r.Count,
	})
}
