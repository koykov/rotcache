package rotcache

import (
	"bytes"
	"context"
	"testing"
	"time"
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
	t.Run("io", func(t *testing.T) {
		c := RotCache{Hasher: testHash{}}
		if err := c.Set("foo", []byte("foobar")); err != nil {
			t.Error("Set failed:", err)
		}
		if err := c.SetString("bar", "foobar"); err != nil {
			t.Error("SetString failed:", err)
		}
		if err := c.USet(10, []byte("foobar")); err != nil {
			t.Error("USet failed:", err)
		}
		if err := c.USetString(10, "foobar"); err != nil {
			t.Error("UsetString failed:", err)
		}
		if err := c.ISet(15, []byte("foobar")); err != nil {
			t.Error("ISet failed:", err)
		}
		if err := c.ISetString(15, "foobar"); err != nil {
			t.Error("ISetString failed:", err)
		}
		c.Rotate()
		if p, err := c.Get("foo"); err != nil || !bytes.Equal(p, []byte("foobar")) {
			t.Error("Get failed:", err)
		}
		if p, err := c.GetString("foo"); err != nil || p != "foobar" {
			t.Error("GetString failed:", err)
		}
		if p, err := c.UGet(10); err != nil || !bytes.Equal(p, []byte("foobar")) {
			t.Error("UGet failed:", err)
		}
		if p, err := c.UGetString(10); err != nil || p != "foobar" {
			t.Error("UGetString failed:", err)
		}
		if p, err := c.IGet(15); err != nil || !bytes.Equal(p, []byte("foobar")) {
			t.Error("IGet failed:", err)
		}
		if p, err := c.IGetString(15); err != nil || p != "foobar" {
			t.Error("IGetString failed:", err)
		}
	})
}

func BenchmarkRotCache(b *testing.B) {
	b.Run("io", func(b *testing.B) {
		c := RotCache{Hasher: testHash{}}
		_ = c.SetString("foo", "qwerty")
		_ = c.SetString("bar", "asdfgh")
		c.Rotate()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if s, err := c.GetString("bar"); err != nil || s != "asdfgh" {
				b.Error("read failed", err)
			}
		}
	})
	b.Run("io parallel", func(b *testing.B) {
		c := RotCache{Hasher: testHash{}}
		write := func(c *RotCache) {
			c.Prepare()
			_ = c.SetString("foo", "qwerty")
			_ = c.SetString("bar", "asdfgh")
			c.Rotate()
		}
		write(&c)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ticker := time.NewTicker(time.Millisecond * 100)
		go func() {
			for {
				select {
				case <-ticker.C:
					write(&c)
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if s, err := c.GetString("bar"); err != nil || s != "asdfgh" {
				b.Error("read failed", err)
			}
		}
		cancel()
	})
}
