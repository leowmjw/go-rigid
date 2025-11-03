package path

import "errors"

type Path struct{ nodes []Navigator }
type Navigator interface{ isNavigator() }

func DecodeJSONPath(data []byte) (Path, error) { return Path{}, errors.New("not implemented") }
func EncodeJSONPath(p Path) ([]byte, error)    { return nil, errors.New("not implemented") }

func Key(k string) Path            { return Path{} }
func Must(p Path) Path             { return Path{} }
func MultiPath(paths ...Path) Path { return Path{} }
func FilterFunc(name string) Path  { return Path{} }
func PKey(v any) Path              { return Path{} }
