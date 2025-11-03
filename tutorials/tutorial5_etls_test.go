//go:build tdd

package tutorials

import (
	"testing"
	"github.com/leowmjw/go-rigid/testutil"
)

func TestT5_Microbatch_ExactlyOnce(t *testing.T) {
	env := testutil.NewEnv()
	defer env.Cancel()
}

func TestT5_Stream_RetryModes_AffectOutputs(t *testing.T) {
}
