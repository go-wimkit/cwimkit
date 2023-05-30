package cwimkit

import (
	"context"

	wasm "github.com/tetratelabs/wazero/api"
)

const AllImages = -1

const (
	OpenCheckIntegrity int32 = 1 << iota
	OpenErrorIfSplit
	OpenWriteAccess
)

const (
	WriteCheckIntegrity int32 = 1 << iota
	WriteNoCheckIntegrity
	WritePipable
	WriteNotPipable
	WriteRecompress
	WriteFsync
	WriteRebuild
	WriteSoftDelete
	WriteIgnoreReadonlyFlag
	WriteSkipExternalWims
	WriteStreamsOk
	WriteRetainGuid
	WriteSolid
	WriteSendDoneWithFileMessages
	WriteNoSolidSort
	WriteUnsafeCompact
)

func (api *API) OpenWIM(ctx context.Context, path string, flags int32) (Pointer, error) {
	outPtr, err := api.malloc(ctx, PointerSize)
	if err != nil {
		return NullPointer, err
	}
	defer func() { _ = api.free(ctx, outPtr) }()

	pathPtr, err := api.dupZString(ctx, &path)
	if err != nil {
		return NullPointer, err
	}
	defer func() { _ = api.free(ctx, pathPtr) }()

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibOpenWim,
		/* const wimlib_tchar *wim_file = */ pathPtr.arg(),
		/* int open_flags = */ wasm.EncodeI32(flags),
		/* WIMStruct **wim_ret = */ outPtr.arg(),
	)
	if err != nil {
		return NullPointer, err
	}

	ref, ok := api.mod.Memory().ReadUint32Le(outPtr.value())
	if !ok {
		return NullPointer, ErrOOBRead
	}

	return Pointer(ref), nil
}

func (api *API) CreateNewWIM(ctx context.Context, ctype CompressionType) (Pointer, error) {
	outPtr, err := api.malloc(ctx, PointerSize)
	if err != nil {
		return NullPointer, err
	}
	defer func() { _ = api.free(ctx, outPtr) }()

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibCreateNewWim,
		/* enum wimlib_compression_type ctype = */ wasm.EncodeI32(int32(ctype)),
		/* WIMStruct **wim_ret = */ outPtr.arg(),
	)
	if err != nil {
		return NullPointer, err
	}

	out, ok := api.mod.Memory().ReadUint32Le(outPtr.value())
	if !ok {
		return NullPointer, ErrOOBRead
	}

	return Pointer(out), nil
}

func (api *API) FreeWIM(ctx context.Context, wim Pointer) error {
	_, err := api.vt.wimlibFree.Call(ctx, wim.arg())
	return err
}

func (api *API) PrintHeader(ctx context.Context, wim Pointer) error {
	_, err := api.vt.wimlibPrintHeader.Call(ctx, wim.arg())
	return err
}

func (api *API) WriteWIM(ctx context.Context, wim Pointer, path string, iid int32, flags int32, threads uint32) error {
	pathPtr, err := api.dupZString(ctx, &path)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, pathPtr) }()

	return api.callWimlibFn(
		ctx,
		api.vt.wimlibWrite,
		/* WIMStruct *wim = */ wim.arg(),
		/* const wimlib_tchar *path = */ pathPtr.arg(),
		/* int image = */ wasm.EncodeI32(iid),
		/* int write_flags */ wasm.EncodeI32(flags),
		/* unsigned num_threads = */ wasm.EncodeU32(threads),
	)
}

func (api *API) OverwriteWIM(ctx context.Context, wim Pointer, flags int32, threads uint32) error {
	return api.callWimlibFn(
		ctx,
		api.vt.wimlibOverwrite,
		/* WIMStruct *wim = */ wim.arg(),
		/* int write_flags = */ wasm.EncodeI32(flags),
		/* unsigned num_threads = */ wasm.EncodeU32(uint32(threads)),
	)
}
