package cwimkit

import (
	"context"
	"errors"
	"sync"

	wasm "github.com/tetratelabs/wazero/api"
)

type eventPortHandler interface {
	Handle(ctx context.Context, inp, out Pointer) error
}

type eventPortHandlerFunc func(ctx context.Context, inp, out Pointer) error

func (f eventPortHandlerFunc) Handle(ctx context.Context, inp, out Pointer) error {
	return f(ctx, inp, out)
}

type eventPortKey struct {
	moduleName   string
	eventPortPtr Pointer
}

type eventPort struct {
	api *API
	ptr Pointer

	handler eventPortHandler
	lastErr error
}

func (p eventPort) free(ctx context.Context) error {
	_, err := p.api.vt.eventPortFree.Call(ctx, p.ptr.arg())
	if err != nil {
		return err
	}

	return nil
}

func (p eventPort) key() eventPortKey {
	return eventPortKey{
		moduleName:   p.api.mod.Name(),
		eventPortPtr: p.ptr,
	}
}

var eventPorts = struct {
	m sync.RWMutex
	h map[eventPortKey]*eventPort
}{h: make(map[eventPortKey]*eventPort)}

func (api *API) eventPortNew(ctx context.Context, h eventPortHandler) (*eventPort, error) {
	ret, err := api.vt.eventPortNew.Call(ctx)
	if err != nil {
		return nil, err
	}

	ep := &eventPort{
		api:     api,
		ptr:     Pointer(ret[0]),
		handler: h,
	}

	eventPorts.m.Lock()
	defer eventPorts.m.Unlock()
	eventPorts.h[ep.key()] = ep

	return ep, nil
}

func processEventFn(ctx context.Context, m wasm.Module, stack []uint64) {
	epPtr := Pointer(wasm.DecodeU32(stack[0]))

	key := eventPortKey{
		moduleName:   m.Name(),
		eventPortPtr: epPtr,
	}

	eventPorts.m.RLock()
	defer eventPorts.m.RUnlock()
	ep, ok := eventPorts.h[key]
	if !ok {
		stack[0] = wasm.EncodeI32(int32(ErrInvalidParam))
		return
	}

	inp := Pointer(wasm.DecodeU32(stack[1]))
	out := Pointer(wasm.DecodeU32(stack[2]))

	var ret int32
	if ep.lastErr = ep.handler.Handle(ctx, inp, out); ep.lastErr != nil {
		var wimlibErr WimlibError
		switch {
		case errors.As(ep.lastErr, &wimlibErr):
			ret = int32(wimlibErr)
		default:
			ret = int32(ErrAbortedByProgress)
		}
	}

	stack[0] = wasm.EncodeI32(ret)
}
