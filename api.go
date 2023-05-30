package cwimkit

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	wasm "github.com/tetratelabs/wazero/api"
	"go.uber.org/multierr"

	"github.com/go-wimkit/cwimkit/internal/csize"
)

type API struct {
	mod wasm.Module
	vt  vtable
}

func (api *API) readZString(offset Pointer) (string, error) {
	var b strings.Builder

	if offset == NullPointer {
		return "", nil
	}

	mem := api.mod.Memory()
	for {
		c, ok := mem.ReadByte(offset.value())
		if !ok {
			return "", ErrOOBRead
		}
		if c == '\x00' {
			return b.String(), nil
		}
		b.WriteByte(c)
		offset += 1
	}
}

func (api *API) readStruct(offset Pointer, value any) error {
	size, ok := csize.Of(value)
	if !ok {
		return fmt.Errorf("readStruct: bad struct")
	}

	buf, ok := api.mod.Memory().Read(offset.value(), size)
	if !ok {
		return ErrOOBRead
	}

	rd := bytes.NewBuffer(buf)
	if err := binary.Read(rd, binary.LittleEndian, value); err != nil {
		return err
	}

	return nil
}

func (api *API) dupZString(ctx context.Context, s *string) (Pointer, error) {
	if s == nil {
		return NullPointer, nil
	}
	n := uint32(len(*s))

	ptr, err := api.malloc(ctx, n+1)
	if err != nil {
		return NullPointer, err
	}

	mem := api.mod.Memory()

	if !mem.WriteString(ptr.value(), *s) {
		return NullPointer, multierr.Append(ErrOOBWrite, api.free(ctx, ptr))
	}

	if !mem.WriteByte(ptr.value()+n, '\x00') {
		return NullPointer, multierr.Append(ErrOOBWrite, api.free(ctx, ptr))
	}

	return ptr, nil
}

func (api *API) dupBytes(ctx context.Context, b []byte) (Pointer, error) {
	ptr, err := api.malloc(ctx, uint32(len(b)))
	if err != nil {
		return NullPointer, err
	}

	if !api.mod.Memory().Write(ptr.value(), b) {
		return NullPointer, multierr.Append(ErrOOBWrite, api.free(ctx, ptr))
	}

	return ptr, nil
}

func (api *API) free(ctx context.Context, addr Pointer) (err error) {
	_, err = api.vt.free.Call(ctx, addr.arg())
	return
}

func (api *API) malloc(ctx context.Context, size uint32) (Pointer, error) {
	ret, err := api.vt.malloc.Call(ctx, wasm.EncodeU32(size))
	if err != nil {
		return NullPointer, err
	}
	if ret[0] == 0 {
		return NullPointer, ErrNoMemory
	}
	return Pointer(ret[0]), err
}

func (api *API) calloc(ctx context.Context, num, size uint32) (Pointer, error) {
	ret, err := api.vt.calloc.Call(ctx, wasm.EncodeU32(num), wasm.EncodeU32(size))
	if err != nil {
		return NullPointer, err
	}
	if ret[0] == 0 {
		return NullPointer, ErrNoMemory
	}
	return Pointer(ret[0]), err
}

func (api *API) fdOpen(ctx context.Context, fd int32, mode string) (Pointer, error) {
	modePtr, err := api.dupZString(ctx, &mode)
	if err != nil {
		return NullPointer, err
	}
	defer func() { _ = api.free(ctx, modePtr) }()

	ret, err := api.vt.fdOpen.Call(ctx, wasm.EncodeI32(fd), modePtr.arg())
	if err != nil {
		return NullPointer, err
	}

	return Pointer(ret[0]), nil
}

func (api *API) callWimlibFn(ctx context.Context, fn wasm.Function, params ...uint64) error {
	ret, err := fn.Call(ctx, params...)
	if err != nil {
		return err
	}

	if rc := ret[0]; rc != 0 {
		return WimlibError(wasm.DecodeI32(rc))
	}

	return nil
}
