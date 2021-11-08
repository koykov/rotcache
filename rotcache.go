package rotcache

import (
	"sync/atomic"

	"github.com/koykov/fastconv"
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
	c.set(c.Hasher.Sum64(key), val)
	return nil
}

func (c *RotCache) SSet(key, val string) error {
	return c.Set(key, fastconv.S2B(val))
}

func (c *RotCache) USet(key uint64, val []byte) error {
	c.set(key, val)
	return nil
}

func (c *RotCache) USSet(key uint64, val string) error {
	c.set(key, fastconv.S2B(val))
	return nil
}

func (c *RotCache) Get(key string) ([]byte, error) {
	if c.Hasher == nil {
		return nil, ErrNoHasher
	}
	return c.get(c.Hasher.Sum64(key)), nil
}

func (c *RotCache) SGet(key string) (string, error) {
	p, err := c.Get(key)
	return fastconv.B2S(p), err
}

func (c *RotCache) UGet(key uint64) ([]byte, error) {
	return c.get(key), nil
}

func (c *RotCache) USGet(key uint64) (string, error) {
	return fastconv.B2S(c.get(key)), nil
}

func (c *RotCache) set(key uint64, val []byte) {
	c.opp().set(key, val)
}

func (c *RotCache) get(key uint64) []byte {
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
