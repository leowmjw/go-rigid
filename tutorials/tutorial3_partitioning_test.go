//go:build tdd

package tutorials

import (
	"testing"

	"example.com/rig/internal/partition"
)

func TestT3_Partitioners_SameKeySamePartition(t *testing.T) {
	h := partition.HashBy(func(k any) int { return 1001 })
	p1, _ := h.Pick("user-123", 64)
	p2, _ := h.Pick("user-123", 64)
	if p1 != p2 { t.Fatalf("expected same partition; %d vs %d", p1, p2) }
}
