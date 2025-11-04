package depot

import (
	"time"

	temporalworkflow "go.temporal.io/sdk/workflow"
)

// AppendRequest is the data sent to the AppendWorkflow.
type AppendRequest struct {
	Module    string
	Depot     string
	Data      any
	AckLevel  AckLevel
}

// AppendWorkflow is the main workflow for appending data to a depot.
func AppendWorkflow(ctx temporalworkflow.Context, req AppendRequest) (AppendResult, error) {
	// This is a simplified implementation. A real implementation would involve
	// a partitioner to determine the target partition, and then would signal
	// or execute a child workflow for that partition.

	// For now, we'll just simulate the different ack levels.
	switch req.AckLevel {
	case None:
		// Fire-and-forget. Return immediately.
		return nil, nil
	case Ack:
		// Acknowledged by the workflow. Return a simple map.
		return AppendResult{"status": "acknowledged"}, nil
	case AppendAck:
		// Acknowledged after processing. We'll simulate this with a timer
		// and then call a placeholder partition workflow.
		var partitionResult any
		ctx = temporalworkflow.WithChildOptions(ctx, temporalworkflow.ChildWorkflowOptions{
			WorkflowExecutionTimeout: 10 * time.Second,
		})
		err := temporalworkflow.ExecuteChildWorkflow(ctx, PartitionWorkflow, req).Get(ctx, &partitionResult)
		if err != nil {
			return nil, err
		}
		return AppendResult{"status": "processed", "partitionResult": partitionResult}, nil
	default:
		return nil, temporalworkflow.NewApplicationError("invalid ack level", "INVALID_ACK_LEVEL")
	}
}

// PartitionWorkflow is a placeholder for the workflow that processes data in a partition.
func PartitionWorkflow(ctx temporalworkflow.Context, req AppendRequest) (any, error) {
	// In a real implementation, this is where the data would be added to a pstate.
	return "partition-ok", nil
}
