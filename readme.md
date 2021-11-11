# Rotating cache

Byte cache implementation with lock-free access. Designed to work with small
data under high pressure. Lock-free access (both read and write) provides by
using two shards of data: one for write, another for read. After finishing
writes cache rotates (swaps) shards.

## Usage

```go
c := rotcache.RotCache{Hasher: fnv.Hasher{}}

go func() {
    for {
        c.Prepare()
        _ = c.SetString("foo", "bar")
        _ = c.ISetString(111, "asd")
        _ = c.SetString("qwe", "rty")
        c.Rotate()
        time.Sleep(time.Second)
        // cancellation logic ...
    }
}()

go func() {
	for {
        v, _ := c.GetString("qwe")
        log.Println(v) // "rty"
        // cancellation logic ...
    }()
}
```
