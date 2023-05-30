package cwimkit

import (
	"context"
	"fmt"

	wasm "github.com/tetratelabs/wazero/api"
)

const (
	IterateDirTreeRecursive       = 0x00000001
	IterateDirTreeChildren        = 0x00000002
	IterateDirTreeResourcesNeeded = 0x00000004
)

type iterateDirTreeEventPortInput struct {
	DirEntry Pointer
}

type iterateDirTreeEventPortHandler struct {
	api *API
	cb  IterateDirTreeCb
}

func (h *iterateDirTreeEventPortHandler) Handle(ctx context.Context, inpPtr, outPtr Pointer) error {
	var inp iterateDirTreeEventPortInput
	if err := h.api.readStruct(inpPtr, &inp); err != nil {
		return err
	}

	dent, err := newDirEntryFromPtr(h.api, inp.DirEntry)
	if err != nil {
		return err
	}

	return h.cb(dent)
}

func (api *API) newIterateDirTreePort(ctx context.Context, cb IterateDirTreeCb) (*eventPort, error) {
	return api.eventPortNew(ctx, &iterateDirTreeEventPortHandler{api: api, cb: cb})
}

type IterateDirTreeCb func(entry *DirEntry) error

func (api *API) IterateDirTree(
	ctx context.Context,
	wim Pointer,
	iid int32,
	path string,
	flags int32,
	cb IterateDirTreeCb,
) error {
	pathPtr, err := api.dupZString(ctx, &path)
	if err != nil {
		return err
	}
	defer func() { _ = api.free(ctx, pathPtr) }()

	p, err := api.newIterateDirTreePort(ctx, cb)
	if err != nil {
		return err
	}
	defer func() { _ = p.free(ctx) }()

	err = api.callWimlibFn(ctx, api.vt.wimlibIterateDirTree, wim.arg(), wasm.EncodeI32(int32(iid)), pathPtr.arg(), wasm.EncodeI32(int32(flags)), p.ptr.arg())
	if err != nil {
		return fmt.Errorf("wimlib.wimlibIterateDirTree: %w", err)
	}

	return nil
}
