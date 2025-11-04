package path

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Path represents a sequence of navigators.
type Path struct {
	nodes []Navigator
}

// Navigator is the interface for all path components.
// The isNavigator method is unexported to prevent external implementations.
type Navigator interface {
	isNavigator()
	json.Marshaler
}

// Navigators is a slice of Navigator, useful for multi-path scenarios.
type Navigators []Navigator

func (n Navigators) isNavigator() {}
func (n Navigators) MarshalJSON() ([]byte, error) {
	items := make([]any, len(n)+1)
	items[0] = "multiPath"
	for i, nav := range n {
		items[i+1] = nav
	}
	return json.Marshal(items)
}

// KeyNav represents a key-based navigation.
type KeyNav struct{ K string }

func (k KeyNav) isNavigator() {}
func (k KeyNav) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.K)
}

// MustNav forces navigation, failing if the path doesn't exist.
type MustNav struct{ Path Path }

func (m MustNav) isNavigator() {}
func (m MustNav) MarshalJSON() ([]byte, error) {
	items := make([]any, len(m.Path.nodes)+1)
	items[0] = "must"
	for i, node := range m.Path.nodes {
		items[i+1] = node
	}
	return json.Marshal(items)
}

// PKeyNav represents a partition key.
type PKeyNav struct{ Val any }

func (p PKeyNav) isNavigator() {}
func (p PKeyNav) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{"pkey", p.Val})
}

// FilterFuncNav represents a filter function.
type FilterFuncNav struct{ Name string }

func (f FilterFuncNav) isNavigator() {}
func (f FilterFuncNav) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Name)
}

// DecodeJSONPath decodes a Path from JSON.
func DecodeJSONPath(data []byte) (Path, error) {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		// If it's not an array of navigators, maybe it's a single navigator?
		// This part of the logic is tricky. Based on the golden files,
		// the top level is always an array.
		return Path{}, err
	}
	return parseNavigators(raw)
}

func parseNavigators(raws []json.RawMessage) (Path, error) {
	var nodes []Navigator
	for _, raw := range raws {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			if strings.HasPrefix(s, "__fOps.") {
				nodes = append(nodes, FilterFuncNav{Name: s})
			} else {
				nodes = append(nodes, KeyNav{K: s})
			}
			continue
		}

		var arr []json.RawMessage
		if err := json.Unmarshal(raw, &arr); err != nil {
			return Path{}, fmt.Errorf("unknown navigator type: %s", string(raw))
		}

		if len(arr) == 0 {
			return Path{}, fmt.Errorf("navigator array cannot be empty")
		}

		var tag string
		if err := json.Unmarshal(arr[0], &tag); err != nil {
			return Path{}, fmt.Errorf("could not unmarshal navigator tag: %s", string(arr[0]))
		}

		switch tag {
		case "pkey":
			if len(arr) != 2 {
				return Path{}, fmt.Errorf("pkey navigator must have 2 elements, got %d", len(arr))
			}
			var val any
			if err := json.Unmarshal(arr[1], &val); err != nil {
				return Path{}, err
			}
			nodes = append(nodes, PKeyNav{Val: val})
		case "must":
			if len(arr) < 2 {
				return Path{}, fmt.Errorf("must navigator must have at least a tag and one path element")
			}
			subPath, err := parseNavigators(arr[1:])
			if err != nil {
				return Path{}, err
			}
			nodes = append(nodes, MustNav{Path: subPath})
		case "multiPath":
			if len(arr) < 2 {
				return Path{}, fmt.Errorf("multiPath navigator must have at least a tag and one path")
			}
			subPath, err := parseNavigators(arr[1:])
			if err != nil {
				return Path{}, err
			}
			nodes = append(nodes, Navigators(subPath.nodes))

		default:
			return Path{}, fmt.Errorf("unknown navigator tag: %s", tag)
		}
	}
	return Path{nodes: nodes}, nil
}

// EncodeJSONPath encodes a Path to JSON.
func EncodeJSONPath(p Path) ([]byte, error) {
	return json.Marshal(p.nodes)
}

// Key creates a new Path starting with a key navigator.
func Key(k string) Path {
	return Path{nodes: []Navigator{KeyNav{K: k}}}
}

// Must creates a new Path starting with a must navigator.
func Must(p Path) Path {
	return Path{nodes: []Navigator{MustNav{Path: p}}}
}

// MultiPath creates a new Path with multiple navigators.
func MultiPath(paths ...Path) Path {
	navs := make([]Navigator, len(paths))
	for i, p := range paths {
		// This is a simplification. A real implementation might flatten the nodes.
		navs[i] = p.nodes[0]
	}
	return Path{nodes: []Navigator{Navigators(navs)}}
}

// FilterFunc creates a new Path with a filter function navigator.
func FilterFunc(name string) Path {
	return Path{nodes: []Navigator{FilterFuncNav{Name: name}}}
}

// PKey creates a new Path with a partition key navigator.
func PKey(v any) Path {
	return Path{nodes: []Navigator{PKeyNav{Val: v}}}
}

// Append adds navigators to an existing path, returning a new Path.
func (p Path) Append(navs ...Navigator) Path {
	newNodes := make([]Navigator, 0, len(p.nodes)+len(navs))
	newNodes = append(newNodes, p.nodes...)
	newNodes = append(newNodes, navs...)
	return Path{nodes: newNodes}
}

// Nodes returns the navigators in the path.
func (p Path) Nodes() []Navigator {
	return p.nodes
}
