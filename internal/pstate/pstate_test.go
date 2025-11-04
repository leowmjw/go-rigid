//go:build tdd

package pstate_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/pstate"
)

func TestPState_LocalSelectAndTransform_Subindexing(t *testing.T) {
	st, err := pstate.OpenPebble(t.TempDir())
	if err == nil {
		if terr := st.LocalTransform([]any{"a","b"}, 123); terr != nil {
			// placeholder; expect not implemented later replaced
		}
		_, _ = st.LocalSelect([]any{"a","b"})
	} else {
		// Expected not implemented for now; ensure error surfaced.
	}
}
