//go:build tdd

package partition_test

import (
	"testing"

	"github.com/leowmjw/go-rigid/internal/partition"
)

func TestPartition_HashBy_Deterministic(t *testing.T) {
	h := partition.HashBy(func(k any) int { return 42 })
	p1, err := h.Pick("user-123", 64)
	if err != nil { t.Fatalf("pick error: %v", err) }
	p2, err := h.Pick("user-123", 64)
	if err != nil { t.Fatalf("pick error: %v", err) }
	if p1 != p2 || p1 != 42%64 {
		t.Fatalf("expected deterministic partition %d got %d,%d", 42%64, p1, p2)
	}
	_, err = h.Pick("user-123", 0)
	if err == nil {
		t.Fatalf("expected error for n=0")
	}
}
