package runner

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

const (
	providerInit     = "PawntestProviderInit"
	providerVersion  = "PawntestProviderVersion"
	providerBefore   = "PawntestProviderBeforeTest"
	providerAfter    = "PawntestProviderAfterTest"
	providerShutdown = "PawntestProviderShutdown"
)

type providerSet struct {
	providers []*pawnProvider
	natives   map[string]*providerNative
}

type pawnProvider struct {
	path     string
	vm       backend.VM
	publics  map[string]backend.Public
	baseline []byte
	active   *providerCall
	registry *providerSet
	maxInstr int
	testName string
}

type providerNative struct {
	provider  *pawnProvider
	handler   string
	signature []providerParamSpec
}

type providerCall struct {
	consumer backend.NativeContext
	params   []backend.Cell
	native   *providerNative
}

type providerParamSpec struct {
	kind        byte
	lengthIndex int
}

func loadProviders(b backend.Backend, paths []string, maxInstructions int) (*providerSet, error) {
	set := &providerSet{natives: map[string]*providerNative{}}

	for _, path := range paths {
		vm, err := b.LoadFile(path)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("load provider %s: %w", path, err), set.Close())
		}

		publics, err := vm.Publics()
		if err != nil {
			_ = vm.Close()
			return nil, errors.Join(fmt.Errorf("read provider %s: %w", path, err), set.Close())
		}

		provider := &pawnProvider{
			path: path, vm: vm, publics: publicMap(publics), registry: set, maxInstr: maxInstructions,
		}

		set.providers = append(set.providers, provider)
		if err := provider.validateVersion(); err != nil {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: %w", path, err), set.Close())
		}

		if _, ok := provider.publics[providerInit]; !ok {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: missing %s", path, providerInit), set.Close())
		}

		if err := provider.registerBridge(); err != nil {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: %w", path, err), set.Close())
		}

		if err := provider.lifecycle(providerInit); err != nil {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: %w", path, err), set.Close())
		}
	}

	return set, nil
}

func (p *pawnProvider) validateVersion() error {
	if _, ok := p.publics[providerVersion]; !ok {
		return errors.New("missing provider ABI version")
	}

	caller, ok := p.vm.(backend.PublicCaller)
	if !ok {
		return errors.New("runtime does not support provider ABI checks")
	}

	version, err := caller.CallPublic(providerVersion)
	if err != nil {
		return err
	}

	if version != 1 {
		return fmt.Errorf("unsupported provider ABI %d; expected 1", version)
	}

	return nil
}

func publicMap(publics []backend.Public) map[string]backend.Public {
	out := make(map[string]backend.Public, len(publics))
	for _, public := range publics {
		out[public.Name] = public
	}

	return out
}

func (p *pawnProvider) registerBridge() error {
	if err := registerFloatNatives(p.vm); err != nil {
		return err
	}

	natives := map[string]backend.NativeFunc{
		"__pawntest_provider_register":   p.registerNative,
		"__pawntest_provider_arg_cell":   p.argCell,
		"__pawntest_provider_arg_string": p.argString,
		"__pt_provider_arg_array":        p.argArrayCell,
		"__pawntest_provider_set_cell":   p.setCell,
		"__pt_provider_set_array":        p.setArrayCell,
		"__pawntest_provider_set_string": p.setString,
		"__pawntest_provider_call":       p.callConsumer,
		"__pt_provider_test_name":        p.writeTestName,
	}
	for name, native := range natives {
		if err := p.vm.RegisterNative(name, native); err != nil {
			return err
		}
	}

	if limiter, ok := p.vm.(backend.InstructionLimiter); ok && p.maxInstr > 0 {
		limiter.SetInstructionLimit(p.maxInstr)
	}

	return nil
}

func (p *pawnProvider) writeTestName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, errors.New("provider test name requires output and size")
	}

	if err := ctx.WriteString(params[0], truncateString(p.testName, int(params[1]))); err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *pawnProvider) registerNative(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 3 {
		return 0, errors.New("provider registration requires native name, handler name, and signature")
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	handler, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	signatureText, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	signature, err := parseProviderSignature(signatureText)
	if err != nil {
		return 0, fmt.Errorf("native %q signature: %w", name, err)
	}

	if name == "" || handler == "" {
		return 0, errors.New("provider registration names cannot be empty")
	}

	if _, ok := p.publics[handler]; !ok {
		return 0, fmt.Errorf("provider handler %q is not public", handler)
	}

	if existing, ok := p.registry.natives[name]; ok {
		return 0, fmt.Errorf("native %q is already provided by %s", name, existing.provider.path)
	}

	p.registry.natives[name] = &providerNative{provider: p, handler: handler, signature: signature}

	return 1, nil
}

func parseProviderSignature(value string) ([]providerParamSpec, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parts := strings.Split(value, ",")

	signature := make([]providerParamSpec, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, errors.New("empty parameter")
		}

		spec := providerParamSpec{kind: part[0], lengthIndex: -1}
		if !strings.ContainsRune("ifsrSAa", rune(spec.kind)) {
			return nil, fmt.Errorf("unsupported parameter kind %q", spec.kind)
		}

		if len(part) > 1 {
			if spec.kind != 'S' && spec.kind != 'A' && spec.kind != 'a' {
				return nil, fmt.Errorf("parameter kind %q cannot have a length", spec.kind)
			}

			if part[1] != ':' {
				return nil, fmt.Errorf("invalid parameter %q", part)
			}

			index, err := strconv.Atoi(part[2:])
			if err != nil || index < 0 || index >= len(parts) {
				return nil, fmt.Errorf("invalid length index in %q", part)
			}

			spec.lengthIndex = index
		}

		if (spec.kind == 'S' || spec.kind == 'A' || spec.kind == 'a') && spec.lengthIndex < 0 {
			return nil, fmt.Errorf("parameter kind %q requires a length index", spec.kind)
		}

		signature = append(signature, spec)
	}

	for _, spec := range signature {
		if spec.lengthIndex >= 0 && signature[spec.lengthIndex].kind != 'i' {
			return nil, fmt.Errorf("length argument %d must have kind i", spec.lengthIndex)
		}
	}

	return signature, nil
}

func (p *pawnProvider) activeCall() (*providerCall, error) {
	if p.active == nil {
		return nil, errors.New("provider argument API called outside a native handler")
	}

	return p.active, nil
}

func providerParam(call *providerCall, index backend.Cell) (backend.Cell, error) {
	if index < 0 || int(index) >= len(call.params) {
		return 0, fmt.Errorf("provider argument index %d is out of range", index)
	}

	return call.params[index], nil
}

func providerSpec(call *providerCall, index backend.Cell, kinds string) (providerParamSpec, error) {
	if _, err := providerParam(call, index); err != nil {
		return providerParamSpec{}, err
	}

	spec := call.native.signature[index]
	if !strings.ContainsRune(kinds, rune(spec.kind)) {
		return providerParamSpec{}, fmt.Errorf("provider argument %d has kind %q; expected %s", index, spec.kind, kinds)
	}

	return spec, nil
}

func providerLength(call *providerCall, spec providerParamSpec) (backend.Cell, error) {
	if spec.lengthIndex < 0 || spec.lengthIndex >= len(call.params) {
		return 0, errors.New("provider parameter has no valid length")
	}

	length := call.params[spec.lengthIndex]
	if length < 0 {
		return 0, errors.New("provider parameter length cannot be negative")
	}

	return length, nil
}

func (p *pawnProvider) argCell(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 1 {
		return 0, errors.New("provider cell input requires an index")
	}

	if _, err := providerSpec(call, params[0], "if"); err != nil {
		return 0, err
	}

	return providerParam(call, params[0])
}

func (p *pawnProvider) argString(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 3 {
		return 0, errors.New("provider string argument requires index, output, and size")
	}

	if _, err := providerSpec(call, params[0], "sS"); err != nil {
		return 0, err
	}

	address, err := providerParam(call, params[0])
	if err != nil {
		return 0, err
	}

	value, err := call.consumer.ReadString(address)
	if err != nil {
		return 0, err
	}

	if err := ctx.WriteString(params[1], truncateString(value, int(params[2]))); err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *pawnProvider) setCell(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 2 {
		return 0, errors.New("provider cell output requires index and value")
	}

	if _, err := providerSpec(call, params[0], "r"); err != nil {
		return 0, err
	}

	address, err := providerParam(call, params[0])
	if err != nil {
		return 0, err
	}

	if err := call.consumer.WriteCell(address, params[1]); err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *pawnProvider) argArrayCell(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 2 {
		return 0, errors.New("provider array input requires index and offset")
	}

	if params[1] < 0 {
		return 0, errors.New("provider array offset cannot be negative")
	}

	spec, err := providerSpec(call, params[0], "aA")
	if err != nil {
		return 0, err
	}

	length, err := providerLength(call, spec)
	if err != nil {
		return 0, err
	}

	if params[1] >= length {
		return 0, fmt.Errorf("provider array offset %d exceeds length %d", params[1], length)
	}

	address, err := providerParam(call, params[0])
	if err != nil {
		return 0, err
	}

	return call.consumer.ReadCell(address + params[1]*4)
}

func (p *pawnProvider) setArrayCell(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 3 {
		return 0, errors.New("provider array output requires index, offset, and value")
	}

	if params[1] < 0 {
		return 0, errors.New("provider array offset cannot be negative")
	}

	spec, err := providerSpec(call, params[0], "A")
	if err != nil {
		return 0, err
	}

	length, err := providerLength(call, spec)
	if err != nil {
		return 0, err
	}

	if params[1] >= length {
		return 0, fmt.Errorf("provider array offset %d exceeds length %d", params[1], length)
	}

	address, err := providerParam(call, params[0])
	if err != nil {
		return 0, err
	}

	if err := call.consumer.WriteCell(address+params[1]*4, params[2]); err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *pawnProvider) setString(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 3 {
		return 0, errors.New("provider string output requires index, value, and size")
	}

	spec, err := providerSpec(call, params[0], "S")
	if err != nil {
		return 0, err
	}

	length, err := providerLength(call, spec)
	if err != nil {
		return 0, err
	}

	if params[2] < 0 || params[2] > length {
		return 0, fmt.Errorf("provider string size %d exceeds length %d", params[2], length)
	}

	address, err := providerParam(call, params[0])
	if err != nil {
		return 0, err
	}

	value, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	if err := call.consumer.WriteString(address, truncateString(value, int(params[2]))); err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *pawnProvider) callConsumer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 1 {
		return 0, errors.New("provider callback requires a public name")
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	caller, ok := call.consumer.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support provider callbacks")
	}

	args := make([]backend.Cell, len(params)-1)
	for i, address := range params[1:] {
		value, err := ctx.ReadCell(address)
		if err != nil {
			return 0, err
		}

		args[i] = value
	}

	return caller.CallPublic(name, args...)
}

func (p *pawnProvider) lifecycle(name string) error {
	if _, ok := p.publics[name]; !ok {
		return nil
	}

	caller, ok := p.vm.(backend.PublicCaller)
	if !ok {
		return errors.New("runtime does not support provider lifecycle callbacks")
	}

	result, err := caller.CallPublic(name)
	if err == nil && result == 0 {
		return fmt.Errorf("%s returned false", name)
	}

	return err
}

func (p *pawnProvider) dispatch(ctx backend.NativeContext, params []backend.Cell, provided *providerNative) (backend.Cell, error) {
	caller, ok := p.vm.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support provider handlers")
	}

	if len(params) != len(provided.signature) {
		return 0, fmt.Errorf("provider native expects %d arguments, received %d", len(provided.signature), len(params))
	}

	previous := p.active

	p.active = &providerCall{consumer: ctx, params: params, native: provided}
	defer func() { p.active = previous }()

	result, err := caller.CallPublic(provided.handler, params...)
	if err != nil {
		return 0, fmt.Errorf("provider %s handler %s: %w", p.path, provided.handler, err)
	}

	return result, nil
}

func (set *providerSet) Register(vm backend.VM, context *executionContext) error {
	for name, provided := range set.natives {
		providerNative := provided

		registered := func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
			if context.mocks.configured(name) {
				return mockUnknownNative(name, context.mocks, context.allowUnknown)(ctx, params)
			}

			context.mocks.recordCall(name, ctx, params)

			return providerNative.provider.dispatch(ctx, params, providerNative)
		}
		if err := vm.RegisterNative(name, registered); err != nil {
			return err
		}
	}

	return nil
}

func (set *providerSet) snapshot() {
	for _, provider := range set.providers {
		if snapshotter, ok := provider.vm.(backend.MemorySnapshotter); ok {
			provider.baseline = snapshotter.SnapshotMemory()
		}
	}
}

func (set *providerSet) beforeTest(isolation, testName string) error {
	var errs []error

	for _, provider := range set.providers {
		provider.testName = testName
		if isolation != isolationSuite && provider.baseline != nil {
			snapshotter, ok := provider.vm.(backend.MemorySnapshotter)
			if !ok {
				errs = append(errs, fmt.Errorf("provider %s: runtime does not support memory snapshots", provider.path))
				continue
			}

			if err := snapshotter.RestoreMemory(provider.baseline); err != nil {
				errs = append(errs, fmt.Errorf("provider %s restore: %w", provider.path, err))
				continue
			}
		}

		if err := provider.lifecycle(providerBefore); err != nil {
			errs = append(errs, fmt.Errorf("provider %s before test: %w", provider.path, err))
		}
	}

	return errors.Join(errs...)
}

func (set *providerSet) afterTest(_ string) error {
	var errs []error

	for _, v := range slices.Backward(set.providers) {
		provider := v
		if err := provider.lifecycle(providerAfter); err != nil {
			errs = append(errs, fmt.Errorf("provider %s after test: %w", provider.path, err))
		}

		provider.testName = ""
	}

	return errors.Join(errs...)
}

func (set *providerSet) Close() error {
	if set == nil {
		return nil
	}

	var errs []error

	for _, v := range slices.Backward(set.providers) {
		provider := v
		errs = append(errs, provider.lifecycle(providerShutdown), provider.vm.Close())
	}

	return errors.Join(errs...)
}
