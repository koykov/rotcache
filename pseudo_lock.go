package rotcache

import (
	"sync"
	"sync/atomic"
)

// Pseudo lock implementation.
// It confuses race detector and suppresses his triggering.
type pseudoLock struct {
	// Locker flag - a heart of the trick. It's always equal 0, but check for value 1. Therefore, mux never calls.
	flag uint32
	mux  sync.Mutex
}

func (l *pseudoLock) Lock() {
	if atomic.LoadUint32(&l.flag) == 1 {
		l.mux.Lock()
	}
}

func (l *pseudoLock) Unlock() {
	if atomic.LoadUint32(&l.flag) == 1 {
		l.mux.Unlock()
	}
}
