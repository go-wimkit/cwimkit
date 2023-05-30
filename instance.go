package cwimkit

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"go.uber.org/multierr"
)

const i32 = wasm.ValueTypeI32

var wimlib struct {
	once      sync.Once
	runtime   wazero.Runtime
	compiled  wazero.CompiledModule
	instances atomic.Uint64
	err       error
}

func compileModule(ctx context.Context) (wazero.Runtime, wazero.CompiledModule, error) {
	r := wazero.NewRuntime(ctx)

	// TODO: close all modules in case of an error.
	_, err := wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		return nil, nil, err
	}

	_, err = r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithGoModuleFunction(
			wasm.GoModuleFunc(processEventFn),
			[]wasm.ValueType{i32, i32, i32},
			[]wasm.ValueType{i32},
		).
		Export("process_event").
		Instantiate(ctx)
	if err != nil {
		return nil, nil, err
	}

	cm, err := r.CompileModule(ctx, assembly.data)
	if err != nil {
		return nil, nil, err
	}

	return r, cm, nil
}

func instantiateModule(ctx context.Context) (wasm.Module, error) {
	wimlib.once.Do(func() {
		wimlib.runtime, wimlib.compiled, wimlib.err = compileModule(ctx)
	})
	if wimlib.err != nil {
		return nil, wimlib.err
	}

	name := "wimlib." + strconv.FormatUint(wimlib.instances.Add(1), 10)

	cfg := wazero.NewModuleConfig().
		WithName(name).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRandSource(rand.Reader).
		WithFSConfig(wazero.NewFSConfig().WithDirMount("/", "/"))

	return wimlib.runtime.InstantiateModule(ctx, wimlib.compiled, cfg)
}

type Instance struct {
	mod wasm.Module
	api *API
}

func (inst *Instance) Close(ctx context.Context) error {
	return inst.mod.Close(ctx)
}

func (inst *Instance) API() *API { return inst.api }

func NewInstance(ctx context.Context) (*Instance, error) {
	mod, err := instantiateModule(ctx)
	if err != nil {
		return nil, fmt.Errorf("instantiateModule: %w", err)
	}

	vt, err := newVTable(mod)
	if err != nil {
		return nil, multierr.Append(err, mod.Close(ctx))
	}

	if _, err := vt.callCtors.Call(ctx); err != nil {
		return nil, multierr.Append(err, mod.Close(ctx))
	}

	inst := &Instance{
		mod: mod,
		api: &API{mod: mod, vt: vt},
	}
	return inst, nil
}
