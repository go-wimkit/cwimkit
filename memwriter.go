package cwimkit

import (
	"math"

	wasm "github.com/tetratelabs/wazero/api"
)

type memoryAtWriter struct {
	mem  wasm.Memory
	base Pointer
}

func (w *memoryAtWriter) WriteAt(p []byte, off int64) (n int, err error) {
	if off > math.MaxUint32 {
		return 0, ErrOOBWrite
	}

	if ok := w.mem.Write(w.base.value()+uint32(off), p); !ok {
		return 0, ErrOOBWrite
	}
	return len(p), nil
}
