package agg

import "fmt"

type Agg interface {
	Fold(state any, value any) (any, error)
}

type CountAgg struct{}

func (CountAgg) Fold(state any, _ any) (any, error) {
	if state == nil {
		return int64(1), nil
	}
	switch s := state.(type) {
	case int:
		return int64(s) + 1, nil
	case int64:
		return s + 1, nil
	default:
		return nil, fmt.Errorf("unexpected state type for Count: %T", state)
	}
}

var Count = CountAgg{}
