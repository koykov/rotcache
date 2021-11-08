package rotcache

import "errors"

var (
	ErrNoHasher    = errors.New("no hasher provided")
	ErrKeyNotFound = errors.New("key not found")
)
