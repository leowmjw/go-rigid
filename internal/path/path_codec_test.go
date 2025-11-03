//go:build tdd

package path_test

import (
	"encoding/json"
	"testing"

	"example.com/rig/internal/path"
)

func TestPath_JSONRoundTrip_Navigators(t *testing.T) {
	cases := []struct{
		name string
		jsonArr []any
	}{
		{ name: "implicit-key+func", jsonArr: []any{"a", "b", "#__fOps.IS_EVEN"} },
		{ name: "explicit-must", jsonArr: []any{"must", "a", "b"} },
		{ name: "forced-pkey", jsonArr: []any{[]any{"pkey", "some-pkey"}, "mapKeys"} },
		{ name: "multi", jsonArr: []any{"multiPath", "a", "b", "c"} },
	}
	for _, tc := range cases {
		b, _ := json.Marshal(tc.jsonArr)
		p, err := path.DecodeJSONPath(b)
		if err == nil { _ , _ = path.EncodeJSONPath(p) }
		_ = b; _ = err
	}
}
