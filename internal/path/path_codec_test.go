//go:build tdd

package path_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leowmjw/go-rigid/internal/path"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestPath_JSONEncoding(t *testing.T) {
	cases := []struct {
		name string
		p    path.Path
	}{
		{
			name: "simple-key",
			p:    path.Key("a"),
		},
		{
			name: "compound-key",
			p:    path.Key("a").Append(path.Key("b").Nodes()...),
		},
		{
			name: "implicit-key-and-func",
			p:    path.Key("a").Append(path.Key("b").Nodes()...).Append(path.FilterFunc("__fOps.IS_EVEN").Nodes()...),
		},
		{
			name: "explicit-must",
			p:    path.Must(path.Key("a").Append(path.Key("b").Nodes()...)),
		},
		{
			name: "forced-pkey",
			p:    path.PKey("some-pkey").Append(path.Key("mapKeys").Nodes()...),
		},
		{
			name: "multi-path",
			p:    path.MultiPath(path.Key("a"), path.Key("b"), path.Key("c")),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			b, err := path.EncodeJSONPath(tc.p)
			require.NoError(t, err)
			g.Assert(t, tc.name, b)
		})
	}
}

func TestPath_JSONRoundTrip_Navigators(t *testing.T) {
	files, err := filepath.Glob("testdata/*.golden")
	require.NoError(t, err)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			golden, err := os.ReadFile(file)
			require.NoError(t, err)

			p, err := path.DecodeJSONPath(golden)
			require.NoError(t, err)

			reencoded, err := path.EncodeJSONPath(p)
			require.NoError(t, err)

			require.JSONEq(t, string(golden), string(reencoded))
		})
	}
}
