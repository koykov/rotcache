package rotcache

import "errors"

var (
	ErrNoHasher    = errors.New("no hasher provided")
	ErrNoKey       = errors.New("no keys provided")
	ErrNoValue     = errors.New("no value provided")
	ErrKeyNotFound = errors.New("key not found")
)
