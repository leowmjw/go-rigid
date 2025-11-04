//go:build tdd

package testutil

import "go.temporal.io/sdk/testsuite"

// WithWorkflowEnv is a small helper to reduce boilerplate in workflow tests.
func WithWorkflowEnv(fn func(env *testsuite.TestWorkflowEnvironment)) {
	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestWorkflowEnvironment()
	fn(env)
}
