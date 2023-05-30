package cwimkit

import (
	wasm "github.com/tetratelabs/wazero/api"
)

type Pointer uint32

const NullPointer Pointer = 0

const PointerSize = 4 // TODO: Oof.

func (p Pointer) arg() uint64 {
	return wasm.EncodeU32(p.value())
}

func (p Pointer) value() uint32 {
	return uint32(p)
}
