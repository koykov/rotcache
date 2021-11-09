package rotcache

import (
	"sync/atomic"

	"github.com/koykov/hash"
)

type RotCache struct {
	Hasher hash.Hasher
	idx    uint32
	buf    [2]cache
}

func (c *RotCache) Set(key string, val []byte) error {
	if c.Hasher == nil {
		return ErrNoHasher
	}
	if len(key) == 0 {
		return ErrNoKey
	}
	if len(val) == 0 {
		return ErrNoValue
	}
	c.set(c.Hasher.Sum64(key), val)
	return nil
}

func (c *RotCache) Get(key string) ([]byte, error) {
	if c.Hasher == nil {
		return nil, ErrNoHasher
	}
	return c.get(c.Hasher.Sum64(key))
}

func (c *RotCache) Prepare() {
	c.opp().reset()
}

func (c *RotCache) Rotate() {
	switch atomic.LoadUint32(&c.idx) {
	case 0:
		atomic.StoreUint32(&c.idx, 1)
	default:
		atomic.StoreUint32(&c.idx, 0)
	}
}

func (c *RotCache) set(key uint64, val []byte) {
	c.opp().set(key, val)
}

func (c *RotCache) get(key uint64) ([]byte, error) {
	return c.act().get(key)
}

func (c *RotCache) act() *cache {
	return (&c.buf[atomic.LoadUint32(&c.idx)]).init()
}

func (c *RotCache) opp() *cache {
	switch atomic.LoadUint32(&c.idx) {
	case 0:
		return (&c.buf[1]).init()
	default:
		return (&c.buf[0]).init()
	}
}
