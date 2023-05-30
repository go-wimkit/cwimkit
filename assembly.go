package cwimkit

import (
	"sync"
)

var assembly struct {
	once sync.Once
	data []byte
}

// SetAssembly sets the wimlib WebAssembly binary to use.
func SetAssembly(data []byte) {
	assembly.once.Do(func() { assembly.data = data })
}
