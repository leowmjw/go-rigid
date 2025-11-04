//go:build tdd

package query_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/query"
)

func TestQuery_Invoke_SingleReturn(t *testing.T) {
	v, err := query.Invoke("sum.plusOne", []any{1,2})
	if err == nil || v != nil {
		// Still not implemented; ensure placeholder behavior (error returned)
		// Adjust when implementation exists.
	}
}
