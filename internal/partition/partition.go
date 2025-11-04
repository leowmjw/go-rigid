package partition

import "errors"

type Partitioner interface{
	Pick(key any, n int) (int, error)
}

type HashBy func(key any) int

func (h HashBy) Pick(key any, n int) (int, error) {
	if n <= 0 {
		return 0, errors.New("invalid partition count")
	}
	hash := h(key)
	return hash % n, nil
}

var Random Partitioner = HashBy(func(any) int { return 0 })
