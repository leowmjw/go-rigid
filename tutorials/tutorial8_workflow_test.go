//go:build tdd

package tutorials

import (
	"testing"

	"github.com/leowmjw/go-rigid/testutil"
	"go.temporal.io/sdk/workflow"
)

// SimpleWorkflow for testing determinism
func SimpleWorkflow(ctx workflow.Context, input string) (string, error) {
	return "processed: " + input, nil
}

func TestT8_WorkflowReplay_Deterministic(t *testing.T) {
	env := testutil.NewEnv()
	// defer env.TearDown()

	env.RegisterWorkflow(SimpleWorkflow)

	env.ExecuteWorkflow(SimpleWorkflow, "test")
	var result1 string
	err := env.GetWorkflowResult(&result1)
	if err != nil {
		t.Fatalf("workflow failed: %v", err)
	}

	env2 := testutil.NewEnv()
	// defer env2.TearDown()
	env2.RegisterWorkflow(SimpleWorkflow)

	env2.ExecuteWorkflow(SimpleWorkflow, "test")
	var result2 string
	err = env2.GetWorkflowResult(&result2)
	if err != nil {
		t.Fatalf("workflow failed: %v", err)
	}

	if result1 != result2 {
		t.Fatalf("workflows not deterministic: %v vs %v", result1, result2)
	}
}
