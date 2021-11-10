package rotcache

import (
	"sync"
)

// Internal cache implementation.
type cache struct {
	once sync.Once
	lock pseudoLock
	// Key-entry pairs storage. Entry value points to offset and length of actual value in buf.
	idx map[uint64]entry
	// Data storage.
	buf []byte
}

func (c *cache) init() *cache {
	// Once apply init logic for each cache instance.
	c.once.Do(func() {
		c.idx = make(map[uint64]entry)
	})
	return c
}

// Set value to the buf using hashed key.
func (c *cache) set(key uint64, val []byte) {
	c.lock.Lock()
	if c.idx == nil {
		c.idx = make(map[uint64]entry)
	}
	var e entry
	e.encode(uln(c.buf), uln(val))
	c.idx[key] = e
	c.buf = append(c.buf, val...)
	c.lock.Unlock()
}

// Get value from the buf using hashed key.
func (c *cache) get(key uint64) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.idx[key]; ok {
		if off, ln := e.decode(); ln < uln(c.buf) {
			return c.buf[off : off+ln], nil
		}
	}
	return nil, ErrKeyNotFound
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
