//go:build tdd

package stream_test

import (
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestStream_EventTree_RetryModes(t *testing.T) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestWorkflowEnvironment()
	defer env.Cancel()
}
