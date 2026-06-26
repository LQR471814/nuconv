package gen

import "testing"

func TestListType(t *testing.T) {
	const expect = `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDeflist = types.List(types.Int())

func NuParselist(v nuplugin.Value) (out list, err error) {
	typed, err := tryCast[[]nuplugin.Value](v)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out = make(list)
	for i, e := range typed {
		out[i], err = nuconvWrapErr(int(e))
		if err != nil {
			err = fmt.Errorf("parse element (%d): %w", i, err)
			return
		}
	}
	return
}
func NuSerializelist(v list) (out nuplugin.Value, err error) {
	tmp := make([]nuplugin.Value, len(v))
	for i, e := range v {
		tmp[i], err = nuconvWrapErr(nuplugin.ToValue(e))
		if err != nil {
			err = fmt.Errorf("serialize element (%d): %w", i, err)
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}
`

	foo := ListType{
		ID: TypeID{
			Pkg:  "testpkg",
			Name: "list",
		},
		ElementType: BuiltinTypeID("int"),
	}

	testType(t, "list.go", foo, expect)
}
