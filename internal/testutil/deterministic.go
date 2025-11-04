//go:build tdd

package testutil

import "time"

// Deterministic stubs for workflows/activities; override in tests if needed.
var NowFunc = func() time.Time { return time.Unix(0, 0).UTC() }
var RandIntn = func(n int) int { if n <= 0 { return 0 }; return 0 }
var UUIDFunc = func() string { return "00000000-0000-0000-0000-000000000000" }
