package depot

import "go.temporal.io/sdk/workflow"

type DepotPartitionArgs struct {
	Module string
	Depot  string
	Part   int
}

func DepotPartitionWorkflow(ctx workflow.Context, args DepotPartitionArgs) error {
	return nil
}
