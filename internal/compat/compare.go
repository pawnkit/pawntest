package compat

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"

	"github.com/pawnkit/pawntest/internal/backend"
)

type PublicResult struct {
	Name  string
	Value backend.Cell
	Error string
}

type Result struct {
	MetadataDiff string
	PublicDiffs  []string
}

func (r Result) Equal() bool {
	return r.MetadataDiff == "" && len(r.PublicDiffs) == 0
}

func Compare(leftName string, left backend.Backend, rightName string, right backend.Backend, tc Case) (Result, error) {
	leftVM, err := left.LoadBytes(tc.Name, tc.AMX)
	if err != nil {
		return Result{}, fmt.Errorf("%s load: %w", leftName, err)
	}
	defer leftVM.Close()
	rightVM, err := right.LoadBytes(tc.Name, tc.AMX)
	if err != nil {
		return Result{}, fmt.Errorf("%s load: %w", rightName, err)
	}
	defer rightVM.Close()

	for name, fn := range tc.Natives {
		if err := leftVM.RegisterNative(name, fn); err != nil {
			return Result{}, fmt.Errorf("%s register native %s: %w", leftName, name, err)
		}
		if err := rightVM.RegisterNative(name, fn); err != nil {
			return Result{}, fmt.Errorf("%s register native %s: %w", rightName, name, err)
		}
	}

	leftMeta, err := metadata(leftVM)
	if err != nil {
		return Result{}, fmt.Errorf("%s metadata: %w", leftName, err)
	}
	rightMeta, err := metadata(rightVM)
	if err != nil {
		return Result{}, fmt.Errorf("%s metadata: %w", rightName, err)
	}
	var result Result
	result.MetadataDiff = cmp.Diff(leftMeta, rightMeta)

	publics, err := leftVM.Publics()
	if err != nil {
		return Result{}, err
	}
	for _, pub := range publics {
		args := tc.PublicArgs[pub.Name]
		leftResult := execPublic(leftVM, pub.Index, pub.Name, args)
		rightResult := execPublic(rightVM, pub.Index, pub.Name, args)
		if diff := cmp.Diff(leftResult, rightResult); diff != "" {
			result.PublicDiffs = append(result.PublicDiffs, diff)
		}
	}
	return result, nil
}

type vmMetadata struct {
	Publics []backend.Public
	Natives []string
}

func metadata(vm backend.VM) (vmMetadata, error) {
	publics, err := vm.Publics()
	if err != nil {
		return vmMetadata{}, err
	}
	natives, err := vm.Natives()
	if err != nil {
		return vmMetadata{}, err
	}
	sort.Slice(publics, func(i, j int) bool { return publics[i].Name < publics[j].Name })
	sort.Strings(natives)
	return vmMetadata{Publics: publics, Natives: natives}, nil
}

func execPublic(vm backend.VM, index int, name string, args []backend.Cell) PublicResult {
	value, err := vm.ExecPublic(index, args...)
	out := PublicResult{Name: name, Value: value}
	if err != nil {
		out.Error = classifyError(err)
	}
	return out
}

func classifyError(err error) string {
	message := strings.ToLower(err.Error())
	for _, item := range []struct {
		kind  string
		terms []string
	}{
		{"divide_by_zero", []string{"divide by zero"}},
		{"bounds", []string{"bounds", "out of bounds"}},
		{"memory", []string{"memory access", "memaccess"}},
		{"instruction", []string{"invalid instruction", "invinstr"}},
		{"sleep", []string{"sleep", "suspended"}},
		{"stack", []string{"stack"}},
		{"heap", []string{"heap"}},
	} {
		for _, term := range item.terms {
			if strings.Contains(message, term) {
				return item.kind
			}
		}
	}
	return message
}
