package cwimkit

import (
	"context"

	wasm "github.com/tetratelabs/wazero/api"
)

const (
	ExportBoot = 1 << iota
	ExportNoNames
	ExportNoDescriptions
	ExportGift
)

func (api *API) ExportImage(
	ctx context.Context, swim Pointer, siid int, dwim Pointer, name, desc *string, flags int,
) (err error) {
	namePtr, err := api.dupZString(ctx, name)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, namePtr) }()

	descPtr, err := api.dupZString(ctx, desc)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, descPtr) }()

	return api.callWimlibFn(
		ctx,
		api.vt.wimlibExportImage,
		/* WIMStruct *src_wim = */ swim.arg(),
		/* int src_image = */ wasm.EncodeI32(int32(siid)),
		/* WIMStruct *dest_wim = */ dwim.arg(),
		/* const wimlib_tchar *dest_name = */ namePtr.arg(),
		/* const wimlib_tchar *dest_description = */ descPtr.arg(),
		/* int export_flags = */ wasm.EncodeI32(int32(flags)),
	)
}

func (api *API) AddEmptyImage(ctx context.Context, wim Pointer, name string) (int, error) {
	outPtr, err := api.malloc(ctx, PointerSize)
	if err != nil {
		return 0, err
	}
	defer func() { _ = api.free(ctx, outPtr) }()

	namePtr, err := api.dupZString(ctx, &name)
	if err != nil {
		return 0, err
	}
	defer func() { _ = api.free(ctx, namePtr) }()

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibAddEmptyImage,
		/* WIMStruct *wim = */ wim.arg(),
		/* wimlib_tchar *name = */ namePtr.arg(),
		/* int *new_idx_ret = */ outPtr.arg(),
	)
	if err != nil {
		return 0, err
	}

	out, ok := api.mod.Memory().ReadUint32Le(outPtr.value())
	if !ok {
		return 0, ErrOOBRead
	}

	return int(out), nil
}

func (api *API) DeleteImage(ctx context.Context, wim Pointer, iid int32) error {
	return api.callWimlibFn(
		ctx,
		api.vt.wimlibDeleteImage,
		/* WIMStruct *wim = */ wim.arg(),
		/* int image = */ wasm.EncodeI32(iid),
	)
}

func (api *API) ResolveImage(ctx context.Context, wim Pointer, nameOrIID string) (int32, error) {
	nameOrIIDPtr, err := api.dupZString(ctx, &nameOrIID)
	if err != nil {
		return 0, err
	}
	defer func() { _ = api.free(ctx, nameOrIIDPtr) }()

	ret, err := api.vt.wimlibResolveImage.Call(
		ctx,
		/* WIMStruct *wim = */ wim.arg(),
		/* const wimlib_tchar *image_name_or_num = */ nameOrIIDPtr.arg(),
	)
	if err != nil {
		return 0, err
	}
	if ret[0] == 0 {
		return 0, ErrInvalidImage
	}
	return int32(ret[0]), nil
}

func (api *API) SetImageProperty(ctx context.Context, wim Pointer, iid int, name string, value *string) error {
	namePtr, err := api.dupZString(ctx, &name)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, namePtr) }()

	valuePtr, err := api.dupZString(ctx, value)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, valuePtr) }()

	return api.callWimlibFn(
		ctx,
		api.vt.wimlibSetImageProperty,
		/* wim = */ wim.arg(),
		/* image = */ wasm.EncodeI32(int32(iid)),
		/* property_name = */ namePtr.arg(),
		/* property_value = */ valuePtr.arg(),
	)
}
