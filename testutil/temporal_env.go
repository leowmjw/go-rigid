package testutil

import "go.temporal.io/sdk/testsuite"

func NewEnv() *testsuite.TestWorkflowEnvironment {
	var ts testsuite.WorkflowTestSuite
	return ts.NewTestWorkflowEnvironment()
}
