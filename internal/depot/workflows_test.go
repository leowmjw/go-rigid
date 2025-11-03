//go:build tdd

package depot_test

import (
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestDepot_Append_AckLevels(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestWorkflowEnvironment()
	defer env.Cancel()
}
