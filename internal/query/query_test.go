//go:build tdd

package query_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/query"
)

func TestQuery_Invoke_SingleReturn(t *testing.T) {
	_, _ = query.Invoke("sum.plusOne", []any{1,2})
}
