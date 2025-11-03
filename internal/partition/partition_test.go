//go:build tdd

package partition_test

import (
	"testing"

	"example.com/rig/internal/partition"
)

func TestPartition_HashBy_Deterministic(t *testing.T) {
	h := partition.HashBy(func(k any) int { return 42 })
	_, _ = h.Pick("user-123", 64)
}
