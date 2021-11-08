package rotcache

import (
	"bytes"
	"testing"
)

type testHash struct{}

func (t testHash) Sum64(s string) uint64 {
	switch s {
	case "foo":
		return 1
	case "bar":
		return 2
	case "qwe":
		return 3
	case "asd":
		return 4
	default:
		return 0
	}
}

func TestRotCache(t *testing.T) {
	t.Run("no hasher", func(t *testing.T) {
		var c RotCache
		err := c.Set("foo", nil)
		if err != ErrNoHasher {
			t.Error("no hasher check failed")
		}
	})
	t.Run("no key", func(t *testing.T) {
		c := RotCache{Hasher: testHash{}}
		err := c.Set("", nil)
		if err != ErrNoKey {
			t.Error("no key check failed")
		}
	})
	t.Run("no value", func(t *testing.T) {
		c := RotCache{Hasher: testHash{}}
		err := c.Set("foo", nil)
		if err != ErrNoValue {
			t.Error("no value check failed")
		}
	})
	t.Run("simple", func(t *testing.T) {
		c := RotCache{Hasher: testHash{}}
		if err := c.Set("foo", []byte("foobar")); err != nil {
			t.Error("set failed:", err)
		}
		if err := c.SSet("bar", "foobar"); err != nil {
			t.Error("sset failed:", err)
		}
		if err := c.USet(10, []byte("foobar")); err != nil {
			t.Error("uset failed:", err)
		}
		if err := c.USSet(10, "foobar"); err != nil {
			t.Error("usset failed:", err)
		}
		c.Rotate()
		if p, err := c.Get("foo"); err != nil || !bytes.Equal(p, []byte("foobar")) {
			t.Error("get failed:", err)
		}
		if p, err := c.SGet("foo"); err != nil || p != "foobar" {
			t.Error("sget failed:", err)
		}
		if p, err := c.UGet(10); err != nil || !bytes.Equal(p, []byte("foobar")) {
			t.Error("uget failed:", err)
		}
		if p, err := c.USGet(10); err != nil || p != "foobar" {
			t.Error("usget failed:", err)
		}
	})
}
