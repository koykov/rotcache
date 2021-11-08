package rotcache

import (
	"sync"

	"github.com/koykov/policy"
)

type cache struct {
	lock policy.Lock
	idx  map[uint64]entry
	buf  []byte
	once sync.Once
}

func (c *cache) init() *cache {
	c.once.Do(func() {
		c.lock.SetPolicy(policy.LockFree)
		c.idx = make(map[uint64]entry)
	})
	return c
}

func (c *cache) set(id uint64, val []byte) {
	c.lock.Lock()
	if c.idx == nil {
		c.idx = make(map[uint64]entry)
	}
	var e entry
	e.encode(uln(c.buf), uln(val))
	c.idx[id] = e
	c.buf = append(c.buf, val...)
	c.lock.Unlock()
}

func (c *cache) get(id uint64) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.idx[id]; ok {
		if off, ln := e.decode(); ln < uln(c.buf) {
			return c.buf[off : off+ln]
		}
	}
	return nil
}

func (c *cache) reset() {
	c.lock.Lock()
	for k := range c.idx {
		delete(c.idx, k)
	}
	c.buf = c.buf[:0]
	c.lock.Unlock()
}

func uln(p []byte) uint32 {
	return uint32(len(p))
}
