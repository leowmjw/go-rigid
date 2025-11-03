package agg

import "errors"

type Agg interface {
	Fold(state any, value any) (any, error)
}

type CountAgg struct{}

func (CountAgg) Fold(state any, _ any) (any, error) { return nil, errors.New("not implemented") }

var Count = CountAgg{}
