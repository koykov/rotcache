package rotcache

import "errors"

var (
	ErrNoHasher = errors.New("no hasher provided")
)
