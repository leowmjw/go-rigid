package partition

import "errors"

type Partitioner interface{
	Pick(key any, n int) (int, error)
}

type HashBy func(key any) int

var ErrInvalidPartitions = errors.New("invalid partitions")

// Pick deterministically maps a key into [0,n). Returns error if n<=0.
func (h HashBy) Pick(key any, n int) (int, error) {
	if n <= 0 {
		return 0, ErrInvalidPartitions
	}
	hv := h(key)
	if hv < 0 {
		hv = -hv
	}
	return hv % n, nil
}

var Random Partitioner = HashBy(func(any) int { return 0 })
