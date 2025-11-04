//go:build tdd

package agg_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/agg"
	"github.com/stretchr/testify/require"
)

func TestAgg_Count(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  any
		iterations    int
		expectedState any
	}{
		{
			name:          "Nil initial state",
			initialState:  nil,
			iterations:    1,
			expectedState: int64(1),
		},
		{
			name:          "Increment once from int",
			initialState:  0,
			iterations:    1,
			expectedState: int64(1),
		},
		{
			name:          "Increment once from int64",
			initialState:  int64(0),
			iterations:    1,
			expectedState: int64(1),
		},
		{
			name:          "Increment three times",
			initialState:  0,
			iterations:    3,
			expectedState: int64(3),
		},
		{
			name:          "Increment many times",
			initialState:  0,
			iterations:    1000,
			expectedState: int64(1000),
		},
		{
			name:          "Increment from non-zero int state",
			initialState:  5,
			iterations:    10,
			expectedState: int64(15),
		},
		{
			name:          "Increment from non-zero int64 state",
			initialState:  int64(5),
			iterations:    10,
			expectedState: int64(15),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.initialState
			var err error
			for i := 0; i < tc.iterations; i++ {
				s, err = agg.Count.Fold(s, 1)
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedState, s)
		})
	}
}