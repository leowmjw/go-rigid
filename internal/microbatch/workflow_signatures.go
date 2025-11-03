package microbatch

import "go.temporal.io/sdk/workflow"

type CoordinatorArgs struct {
	Module   string
	Topology string
	Task     int
}

func CoordinatorWorkflow(ctx workflow.Context, args CoordinatorArgs) error {
	return ErrNotImplemented
}

type AttemptArgs struct {
	Module   string
	Topology string
	BatchID  string
}

func AttemptWorkflow(ctx workflow.Context, args AttemptArgs) error {
	return ErrNotImplemented
}
