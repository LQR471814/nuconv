package logic

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/require"
)

func testType(t *testing.T, typ Type, expected string) {
	c := GenContext{
		Out: jen.NewFilePath("test/pkg"),
	}

	typ.Definition(c)
	typ.Parser(c)
	typ.Serializer(c)

	buf := bytes.NewBuffer(nil)
	err := c.Out.Render(buf)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, buf.String())
}

func TestArrayType(t *testing.T) {
	const expect = `package pkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDefFoo = types.List(string)

func NuParseFoo(v nuplugin.Value) (out Foo, err error) {
	typed, ok := v.Value.([]nuplugin.Value)
	if !ok {
		err = fmt.Errorf("expected %T got %T", typed, v)
		return
	}
	if len(typed) != 3 {
		err = fmt.Errorf("array has wrong number of elements, got (%d), expected (%d)", len(typed), 3)
		return
	}
	out = make(Foo)
	for i := range 3 {
		out[i], err = nuconvWrapErr(string(typed[i]))
		if err != nil {
			return
		}
	}
	return
}
func NuSerializeFoo(v Foo) (out nuplugin.Value, err error) {
	tmp := make([]nuplugin.Value, 3)
	for i := range 3 {
		tmp[i], err = nuconvWrapErr(nuplugin.ToValue(typed[i]))
		if err != nil {
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}
`

	foo := ArrayType{
		ID: TypeID{
			Pkg:  "test/pkg",
			Name: "Foo",
		},
		Length:      3,
		ElementType: BuiltinTypeID("string"),
	}

	testType(t, foo, expect)
}
