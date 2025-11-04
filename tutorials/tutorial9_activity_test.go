//go:build tdd

package tutorials

import (
	"context"
	"errors"
	"testing"

	"github.com/leowmjw/go-rigid/testutil"
	"go.temporal.io/sdk/workflow"
)

// FailingActivity for testing retries
func FailingActivity(ctx context.Context) error {
	return errors.New("activity failed")
}

func WorkflowWithActivity(ctx workflow.Context) error {
	return workflow.ExecuteActivity(ctx, FailingActivity).Get(ctx, nil)
}

func TestT9_ActivityRetries_OnFailure(t *testing.T) {
	env := testutil.NewEnv()
	// defer env.TearDown()

	env.RegisterWorkflow(WorkflowWithActivity)
	env.RegisterActivity(FailingActivity)

	env.ExecuteWorkflow(WorkflowWithActivity)
	err := env.GetWorkflowResult(nil)
	if err == nil {
		t.Fatalf("expected workflow to fail")
	}
	// TODO: assert retry attempts
}
