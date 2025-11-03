//go:build tdd

package agg_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/agg"
)

func TestAgg_Count(t *testing.T) {
	var s any
	for i:=0; i<3; i++ {
		var err error
		s, err = agg.Count.Fold(s, 1)
		if err != nil { t.Fatalf("fold: %v", err) }
	}
}
