package rotcache

// Extended setters and getters.

import "github.com/koykov/fastconv"

// SetString sets string value with string key.
func (c *RotCache) SetString(key, val string) error {
	return c.Set(key, fastconv.S2B(val))
}

// USet sets new value with uint64 key.
func (c *RotCache) USet(key uint64, val []byte) error {
	if len(val) == 0 {
		return ErrNoValue
	}
	c.set(key, val)
	return nil
}

// USetString sets string value with uint64 key.
func (c *RotCache) USetString(key uint64, val string) error {
	if len(val) == 0 {
		return ErrNoValue
	}
	c.set(key, fastconv.S2B(val))
	return nil
}

// ISet sets new value with int64 key.
func (c *RotCache) ISet(key int64, val []byte) error {
	return c.USet(uint64(key), val)
}

// ISetString sets string value with int64 key.
func (c *RotCache) ISetString(key int64, val string) error {
	return c.USetString(uint64(key), val)
}

// GetString gets stored value as string.
func (c *RotCache) GetString(key string) (string, error) {
	p, err := c.Get(key)
	return fastconv.B2S(p), err
}

// UGet get stored value using uint64 key.
func (c *RotCache) UGet(key uint64) ([]byte, error) {
	return c.get(key)
}

// UGetString get stored value as string using uint64 key.
func (c *RotCache) UGetString(key uint64) (s string, err error) {
	var p []byte
	if p, err = c.get(key); err != nil {
		return "", err
	}
	return fastconv.B2S(p), err
}

// IGet get stored value using int64 key.
func (c *RotCache) IGet(key int64) ([]byte, error) {
	return c.get(uint64(key))
}

// IGetString get stored value as string using int64 key.
func (c *RotCache) IGetString(key int64) (s string, err error) {
	return c.UGetString(uint64(key))
}
