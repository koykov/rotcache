package rotcache

import "sync/atomic"

// Internal cache implementation.
type cache struct {
	// Locking flag.
	lock uint32
	// Key-entry pairs storage. Entry value points to offset and length of actual value in buf.
	idx []idx
	// Data storage.
	buf []byte
}

type idx struct {
	h uint64
	e entry
}

// Set value to the buf using hashed key.
//go:norace
func (c *cache) set(key uint64, val []byte) {
	c.waitLock()
	atomic.StoreUint32(&c.lock, 1)
	defer atomic.StoreUint32(&c.lock, 0)
	var e entry
	e.encode(uln(c.buf), uln(val))
	c.idx = append(c.idx, idx{
		h: key,
		e: e,
	})
	c.buf = append(c.buf, val...)
}

// Get value from the buf using hashed key.
func (c *cache) get(key uint64) ([]byte, error) {
	c.waitLock()
	l := len(c.idx)
	if l == 0 {
		return nil, ErrKeyNotFound
	}
	_ = c.idx[l-1]
	for i := 0; i < l; i++ {
		if c.idx[i].h == key {
			e := c.idx[i].e
			if off, ln := e.decode(); ln < uln(c.buf) {
				return c.buf[off : off+ln], nil
			}
		}
	}
	return nil, ErrKeyNotFound
}

//go:norace
func (c *cache) reset() {
	c.waitLock()
	c.idx = c.idx[:0]
	c.buf = c.buf[:0]
}

func (c *cache) waitLock() {
	for atomic.LoadUint32(&c.lock) == 1 {
	}
}

func uln(p []byte) uint32 {
	return uint32(len(p))
}
