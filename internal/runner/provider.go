package runner

import (
	"errors"
	"fmt"
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

const (
	providerInit     = "PawntestProviderInit"
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
}

type providerNative struct {
	provider *pawnProvider
	handler  string
}

type providerCall struct {
	consumer backend.NativeContext
	params   []backend.Cell
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

		if err := provider.registerBridge(); err != nil {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: %w", path, err), set.Close())
		}

		if err := provider.lifecycle(providerInit); err != nil {
			return nil, errors.Join(fmt.Errorf("initialize provider %s: %w", path, err), set.Close())
		}
	}

	return set, nil
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

func (p *pawnProvider) registerNative(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, errors.New("provider registration requires native and handler names")
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	handler, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
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

	p.registry.natives[name] = &providerNative{provider: p, handler: handler}

	return 1, nil
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

func (p *pawnProvider) argCell(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	call, err := p.activeCall()
	if err != nil {
		return 0, err
	}

	if len(params) < 1 {
		return 0, errors.New("provider cell input requires an index")
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

	_, err := caller.CallPublic(name)

	return err
}

func (p *pawnProvider) dispatch(ctx backend.NativeContext, params []backend.Cell, handler string) (backend.Cell, error) {
	caller, ok := p.vm.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support provider handlers")
	}

	previous := p.active

	p.active = &providerCall{consumer: ctx, params: params}
	defer func() { p.active = previous }()

	result, err := caller.CallPublic(handler, params...)
	if err != nil {
		return 0, fmt.Errorf("provider %s handler %s: %w", p.path, handler, err)
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

			return providerNative.provider.dispatch(ctx, params, providerNative.handler)
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

func (set *providerSet) beforeTest(isolation, _ string) error {
	for _, provider := range set.providers {
		if isolation != isolationSuite && provider.baseline != nil {
			snapshotter, ok := provider.vm.(backend.MemorySnapshotter)
			if !ok {
				return errors.New("provider runtime does not support memory snapshots")
			}

			if err := snapshotter.RestoreMemory(provider.baseline); err != nil {
				return err
			}
		}

		if err := provider.lifecycle(providerBefore); err != nil {
			return fmt.Errorf("provider %s before test: %w", provider.path, err)
		}
	}

	return nil
}

func (set *providerSet) afterTest(_ string) error {
	for _, v := range slices.Backward(set.providers) {
		provider := v
		if err := provider.lifecycle(providerAfter); err != nil {
			return fmt.Errorf("provider %s after test: %w", provider.path, err)
		}
	}

	return nil
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
