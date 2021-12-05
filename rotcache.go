package rotcache

import (
	"sync/atomic"

	"github.com/koykov/hash"
)

// RotCache is a rotating cache implementation.
type RotCache struct {
	// Keys hasher. Uses to convert string keys to hashes.
	Hasher hash.Hasher
	// Index of currently uses cache.
	// May only take two values:
	// * 0 - use buf[0]
	// * 1 - use buf[1]
	idx uint32
	// Two internal caches.
	// Only one can be used at the same time, dependent of idx.
	// The cache used currently calls "actual", otherwise calls "opposite".
	buf [2]cache
}

// Set sets new key-value pair to the opposite cache.
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

// Get gets stored value by key from the actual cache.
func (c *RotCache) Get(key string) ([]byte, error) {
	if c.Hasher == nil {
		return nil, ErrNoHasher
	}
	return c.get(c.Hasher.Sum64(key))
}

// Prepare resets the opposite cache.
func (c *RotCache) Prepare() {
	c.opp().reset()
}

// Rotate swaps actual and opposite caches.
func (c *RotCache) Rotate() {
	switch atomic.LoadUint32(&c.idx) {
	case 0:
		atomic.StoreUint32(&c.idx, 1)
	default:
		atomic.StoreUint32(&c.idx, 0)
	}
}

// Internal setter method.
func (c *RotCache) set(key uint64, val []byte) {
	c.opp().set(key, val)
}

// Internal getter method.
func (c *RotCache) get(key uint64) ([]byte, error) {
	return c.act().get(key)
}

// Get instance of actual cache.
func (c *RotCache) act() *cache {
	return &c.buf[atomic.LoadUint32(&c.idx)]
}

// Get instance of opposite cache.
func (c *RotCache) opp() *cache {
	switch atomic.LoadUint32(&c.idx) {
	case 0:
		return &c.buf[1]
	default:
		return &c.buf[0]
	}
}
