package cwimkit

import (
	"context"
	"io"

	wasm "github.com/tetratelabs/wazero/api"
)

type UpdateOp int32

const (
	UpdateOpAdd UpdateOp = iota
	UpdateOpDelete
	UpdateOpRename
)

const (
	cmdBufItemSize = 20
)

type ResourceOwner interface {
	Free(ctx context.Context, api *API) error
}

type WriterToProvider interface {
	WriterTo(ctx context.Context, api *API) (io.WriterTo, error)
}

func (api *API) UpdateImage(
	ctx context.Context,
	wim Pointer,
	iid int32,
	flags int32,
	cmds ...WriterToProvider,
) error {
	if len(cmds) == 0 {
		return nil
	}

	resources := make([]io.WriterTo, len(cmds))
	defer func() {
		for _, res := range resources {
			if res == nil {
				continue
			}
			if o, ok := res.(ResourceOwner); ok {
				_ = o.Free(ctx, api)
			}
		}
	}()

	for i, cmd := range cmds {
		var err error
		resources[i], err = cmd.WriterTo(ctx, api)
		if err != nil {
			return err
		}
	}

	cmdBufPtr, err := api.calloc(ctx, uint32(len(cmds)), cmdBufItemSize)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, cmdBufPtr) }()

	memWr := &memoryAtWriter{
		mem:  api.mod.Memory(),
		base: cmdBufPtr,
	}
	for i, res := range resources {
		wr := io.NewOffsetWriter(memWr, int64(cmdBufItemSize*i))
		if _, err := res.WriteTo(wr); err != nil {
			return err
		}
	}

	err = api.callWimlibFn(
		ctx,
		api.vt.wimlibUpdateImage,
		/* WIMStruct *wim = */ wim.arg(),
		/* int image = */ wasm.EncodeI32(iid),
		/* const struct wimlib_update_command *cmds = */ cmdBufPtr.arg(),
		/* size_t num_cmds = */ wasm.EncodeU32(uint32(len(cmds))),
		/* int update_flags = */ wasm.EncodeI32(flags),
	)
	if err != nil {
		return err
	}

	return nil
}
