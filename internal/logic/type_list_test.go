package logic

import "testing"

func TestListType(t *testing.T) {
	const expect = `package pkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDeffoo = types.List(int)

func NuParsefoo(v nuplugin.Value) (out foo, err error) {
	typed, err := tryCast[[]nuplugin.Value](v)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out = make(foo)
	for i, e := range typed {
		out[i], err = nuconvWrapErr(int(e))
		if err != nil {
			err = fmt.Errorf("parse element (%d): %w", i, err)
			return
		}
	}
	return
}
func NuSerializefoo(v foo) (out nuplugin.Value, err error) {
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
			Pkg:  "test/pkg",
			Name: "foo",
		},
		ElementType: BuiltinTypeID("int"),
	}

	testType(t, foo, expect)
}
