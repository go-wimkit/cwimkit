package cwimkit

import (
	"context"

	wasm "github.com/tetratelabs/wazero/api"
)

type CompressionType int32

const (
	CompressionTypeXPRESS CompressionType = iota + 1
	CompressionTypeLZX
	CompressionTypeLZMS
)

func (api *API) CreateDecompressor(ctx context.Context, ctype, chunkLen int32) (Pointer, error) {
	ptr, err := api.malloc(ctx, 4)
	if err != nil {
		return NullPointer, err
	}
	defer func() { _ = api.free(ctx, ptr) }()

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibCreateDecompressor,
		/* ctype = */ wasm.EncodeI32(ctype),
		/* max_block_size = */ wasm.EncodeI32(chunkLen),
		/* decompressor_ret = */ ptr.arg(),
	)
	if err != nil {
		return NullPointer, err
	}

	dec, ok := api.mod.Memory().ReadUint32Le(ptr.value())
	if !ok {
		return NullPointer, ErrOOBRead
	}

	return Pointer(dec), nil
}

func (api *API) Decompress(ctx context.Context, dec Pointer, inp, out []byte) error {
	inpPtr, err := api.dupBytes(ctx, inp)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, inpPtr) }()

	outPtr, err := api.malloc(ctx, uint32(len(out)))
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, outPtr) }()

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibDecompress,
		/* compressed_data = */ inpPtr.arg(),
		/* compressed_size = */ wasm.EncodeU32(uint32(len(inp))),
		/* uncompressed_data = */ outPtr.arg(),
		/* uncompressed_size = */ wasm.EncodeU32(uint32(len(out))),
		/* dec = */ dec.arg(),
	)
	if err != nil {
		return err
	}

	mem, ok := api.mod.Memory().Read(outPtr.value(), uint32(len(out)))
	if !ok {
		return ErrOOBRead
	}
	copy(out, mem)

	return nil
}

func (api *API) FreeDecompressor(ctx context.Context, dec Pointer) (err error) {
	_, err = api.vt.wimlibFreeDecompressor.Call(ctx, dec.arg())
	return
}
