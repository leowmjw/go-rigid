//go:build tdd

package pstate_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/pstate"
)

func TestPState_LocalSelectAndTransform_Subindexing(t *testing.T) {
	st, err := pstate.OpenPebble(t.TempDir())
	if err == nil {
		_ = st.LocalTransform([]any{"a","b"}, 123)
		_, _ = st.LocalSelect([]any{"a","b"})
	}
}
