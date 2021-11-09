package rotcache

import (
	"sync"
	"sync/atomic"
)

type pseudoLock struct {
	flag uint32
	mux sync.Mutex
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
