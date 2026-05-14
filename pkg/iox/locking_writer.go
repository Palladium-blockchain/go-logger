package iox

import (
	"sync"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type LockingWriter struct {
	w  core.Writer
	mu sync.Mutex
}

func NewLockingWriter(w core.Writer) *LockingWriter {
	return &LockingWriter{w: w}
}

func (w *LockingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.w.Write(p)
}
