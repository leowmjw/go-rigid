package depot

import "context"

type AckLevel string

const (
	Ack AckLevel = "ack"
	AppendAck AckLevel = "appendAck"
	None AckLevel = "none"
)

type AppendResult map[string]any

func ClientAppend(ctx context.Context, module, depot string, data any, ack AckLevel) (AppendResult, error) {
	return nil, ErrNotImplemented
}
