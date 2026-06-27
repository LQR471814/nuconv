package gen

import "testing"

func TestTypeOptional(t *testing.T) {
	expected := `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDefOptional = types.OneOf(types.Nothing(), types.String())

func NuParseOptional(v nuplugin.Value) (out Optional, err error) {
	if v.Value == nil {
		return
	}
	var tmp string
	tmp, err = tryCast[string](v.Value)
	if err != nil {
		err = fmt.Errorf("parse underlying: %w", err)
		return
	}
	out = &tmp
	return
}
func NuSerializeOptional(v Optional) (out nuplugin.Value, err error) {
	if v == nil {
		out = nuplugin.ToValue(nil)
		return
	}
	out, err = nuconvWrapErr(nuplugin.ToValue(*v))
	if err != nil {
		err = fmt.Errorf("serialize underlying: %w", err)
		return
	}
	return
}
`

	optional := OptionalType{
		ID: TypeID{
			Pkg:  "testpkg",
			Name: "Optional",
		},
		Elem: BuiltinTypeID("string"),
	}

	testType(t, "optional.go", optional, expected)
}
