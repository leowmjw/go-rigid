//go:build tdd

package tutorials

import (
	"testing"
	"github.com/leowmjw/go-rigid/testutil"
)

func TestT5_Microbatch_ExactlyOnce(t *testing.T) {
	_ = testutil.NewEnv()
	// TODO: implement microbatch exactly-once test
}

func TestT5_Stream_RetryModes_AffectOutputs(t *testing.T) {
}
