package stream

import "go.temporal.io/sdk/workflow"

type ProcessorArgs struct {
	Module   string
	Topology string
	Part     int
}

func StreamProcessorWorkflow(ctx workflow.Context, args ProcessorArgs) error {
	return ErrNotImplemented
}
