package runner

import (
	"errors"
	"strings"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestParseProviderSignature(t *testing.T) {
	signature, err := parseProviderSignature("i,S:2,i,a:4,i")
	if err != nil {
		t.Fatal(err)
	}

	if len(signature) != 5 || signature[1].kind != 'S' || signature[1].lengthIndex != 2 || signature[3].lengthIndex != 4 {
		t.Fatalf("signature = %#v", signature)
	}
}

func FuzzParseProviderSignature(f *testing.F) {
	f.Add("i,S:2,i")
	f.Add("")
	f.Add("a:999")
	f.Fuzz(func(t *testing.T, signature string) {
		_, _ = parseProviderSignature(signature)
	})
}

func TestParseProviderSignatureRejectsInvalidParameters(t *testing.T) {
	for _, value := range []string{"x", "S", "i:1", "a:9", "i,,i"} {
		if _, err := parseProviderSignature(value); err == nil {
			t.Errorf("signature %q was accepted", value)
		}
	}
}

func TestProviderBridgeValidatesKindsAndBounds(t *testing.T) {
	provider := &pawnProvider{active: &providerCall{
		consumer: newFakeVM(),
		params:   []backend.Cell{100, 2},
		native: &providerNative{signature: []providerParamSpec{
			{kind: 'A', lengthIndex: 1},
			{kind: 'i', lengthIndex: -1},
		}},
	}}

	if _, err := provider.setArrayCell(nil, []backend.Cell{0, 2, 10}); err == nil || !strings.Contains(err.Error(), "exceeds length") {
		t.Fatalf("array bounds error = %v", err)
	}

	if _, err := provider.setCell(nil, []backend.Cell{0, 10}); err == nil || !strings.Contains(err.Error(), "expected r") {
		t.Fatalf("reference kind error = %v", err)
	}
}

func TestProviderRegistrationRejectsDuplicates(t *testing.T) {
	context := newFakeVM()
	context.strings[100] = "Inventory_Add"
	context.strings[200] = "Provider_Add"
	context.strings[300] = "i"
	set := &providerSet{natives: map[string]*providerNative{}}
	provider := &pawnProvider{path: "inventory.amx", publics: map[string]backend.Public{"Provider_Add": {}}, registry: set}

	if _, err := provider.registerNative(context, []backend.Cell{100, 200, 300}); err != nil {
		t.Fatal(err)
	}

	if _, err := provider.registerNative(context, []backend.Cell{100, 200, 300}); err == nil || !strings.Contains(err.Error(), "already provided") {
		t.Fatalf("duplicate error = %v", err)
	}
}

func TestProviderVersionAndLifecycleResults(t *testing.T) {
	vm := &providerTestVM{fakeVM: newFakeVM(), results: map[string]backend.Cell{providerVersion: 1, providerBefore: 0}}
	provider := &pawnProvider{vm: vm, publics: map[string]backend.Public{providerVersion: {}, providerBefore: {}}}

	if err := provider.validateVersion(); err != nil {
		t.Fatal(err)
	}

	if err := provider.lifecycle(providerBefore); err == nil || !strings.Contains(err.Error(), "returned false") {
		t.Fatalf("lifecycle error = %v", err)
	}

	vm.results[providerVersion] = 2

	if err := provider.validateVersion(); err == nil || !strings.Contains(err.Error(), "unsupported provider ABI") {
		t.Fatalf("version error = %v", err)
	}
}

type providerTestVM struct {
	*fakeVM
	results map[string]backend.Cell
	err     error
}

func (vm *providerTestVM) CallPublic(name string, args ...backend.Cell) (backend.Cell, error) {
	if vm.err != nil {
		return 0, vm.err
	}

	result, ok := vm.results[name]
	if !ok {
		return 0, errors.New("missing public")
	}

	return result, nil
}
