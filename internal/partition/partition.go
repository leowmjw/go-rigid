package partition

import "errors"

type Partitioner interface{
	Pick(key any, n int) (int, error)
}

type HashBy func(key any) int

func (h HashBy) Pick(key any, n int) (int, error) { return 0, errors.New("not implemented") }

var Random Partitioner = HashBy(func(any) int { return 0 })
