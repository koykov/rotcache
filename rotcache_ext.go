package rotcache

import "github.com/koykov/fastconv"

func (c *RotCache) SetString(key, val string) error {
	return c.Set(key, fastconv.S2B(val))
}

func (c *RotCache) USet(key uint64, val []byte) error {
	if len(val) == 0 {
		return ErrNoValue
	}
	c.set(key, val)
	return nil
}

func (c *RotCache) USetString(key uint64, val string) error {
	if len(val) == 0 {
		return ErrNoValue
	}
	c.set(key, fastconv.S2B(val))
	return nil
}

func (c *RotCache) ISet(key int64, val []byte) error {
	return c.USet(uint64(key), val)
}

func (c *RotCache) ISetString(key int64, val string) error {
	return c.USetString(uint64(key), val)
}

func (c *RotCache) GetString(key string) (string, error) {
	p, err := c.Get(key)
	return fastconv.B2S(p), err
}

func (c *RotCache) UGet(key uint64) ([]byte, error) {
	return c.get(key)
}

func (c *RotCache) UGetString(key uint64) (s string, err error) {
	var p []byte
	if p, err = c.get(key); err != nil {
		return "", err
	}
	return fastconv.B2S(p), err
}

func (c *RotCache) IGet(key int64) ([]byte, error) {
	return c.get(uint64(key))
}

func (c *RotCache) IGetString(key int64) (s string, err error) {
	return c.UGetString(uint64(key))
}
