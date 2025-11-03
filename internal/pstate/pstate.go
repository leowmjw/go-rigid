package pstate

import "errors"

type Store interface {
	LocalSelect(path any) (any, error)
	LocalTransform(path any, value any) error
}

type KV interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Batch(fn func(KV) error) error
}

var ErrNotImplemented = errors.New("not implemented")

type pebbleStore struct{}

func OpenPebble(path string) (Store, error) { return nil, ErrNotImplemented }
