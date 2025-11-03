//go:build tdd

package microbatch_test

import (
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestMicrobatch_ExactlyOnce_AcrossAttempts(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestWorkflowEnvironment()
	defer env.Cancel()
}
