package depot

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

type AckLevel string

const (
	Ack       AckLevel = "ack"
	AppendAck AckLevel = "appendAck"
	None      AckLevel = "none"
)

// AppendResult is the result of an append operation.
type AppendResult map[string]any

// ClientAppend starts the AppendWorkflow.
func ClientAppend(ctx context.Context, c client.Client, module, depot string, data any, ack AckLevel) (AppendResult, error) {
	req := AppendRequest{
		Module:   module,
		Depot:    depot,
		Data:     data,
		AckLevel: ack,
	}

	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("depot-append-%s-%s", module, depot), // This should be more unique in a real app
		TaskQueue: "depot-tasks",
	}

	wfRun, err := c.ExecuteWorkflow(ctx, options, AppendWorkflow, req)
	if err != nil {
		return nil, err
	}

	if ack == None {
		return nil, nil
	}

	var result AppendResult
	if err := wfRun.Get(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
