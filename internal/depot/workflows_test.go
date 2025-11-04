//go:build tdd

package depot_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/depot"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

func TestDepot_Append_AckLevels(t *testing.T) {
	testCases := []struct {
		name         string
		ackLevel     depot.AckLevel
		expectResult depot.AppendResult
		expectError  bool
	}{
		{
			name:         "AckLevel None",
			ackLevel:     depot.None,
			expectResult: nil,
			expectError:  false,
		},
		{
			name:         "AckLevel Ack",
			ackLevel:     depot.Ack,
			expectResult: depot.AppendResult{"status": "acknowledged"},
			expectError:  false,
		},
		{
			name:         "AckLevel AppendAck",
			ackLevel:     depot.AppendAck,
			expectResult: depot.AppendResult{"status": "processed", "partitionResult": "partition-ok"},
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ts testsuite.WorkflowTestSuite
			env := ts.NewTestWorkflowEnvironment()
			defer env.Cancel()

			env.RegisterWorkflow(depot.AppendWorkflow)
			env.RegisterWorkflow(depot.PartitionWorkflow)

			req := depot.AppendRequest{
				Module:   "test-module",
				Depot:    "test-depot",
				Data:     "test-data",
				AckLevel: tc.ackLevel,
			}

			env.ExecuteWorkflow(depot.AppendWorkflow, req)

			require.True(t, env.IsWorkflowCompleted())
			if tc.expectError {
				require.Error(t, env.GetWorkflowError())
			} else {
				require.NoError(t, env.GetWorkflowError())
				var result depot.AppendResult
				env.GetWorkflowResult(&result)
				require.Equal(t, tc.expectResult, result)
			}
		})
	}
}
