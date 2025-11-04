package agg

import "fmt"

type Agg interface {
	Fold(state any, value any) (any, error)
}

type CountAgg struct{}

// Fold increments an integer counter; nil state initializes to 0.
func (CountAgg) Fold(state any, _ any) (any, error) {
	if state == nil {
		return 1, nil
	}
	v, ok := state.(int)
	if !ok {
		return nil, fmt.Errorf("count state type %T not int", state)
	}
	return v + 1, nil
}

var Count = CountAgg{}
