package gen

import "testing"

func TestTypeMap(t *testing.T) {
	expected := `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDefMap = types.Record(types.RecordDef{})

func NuParseMap(v nuplugin.Value) (out Map, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out = make(map[string]int)
	for key, e := range typed {
		out[key], err = tryCast[int](e.Value)
		if err != nil {
			err = fmt.Errorf("parse map entry (%s:%v): %w", key, e, err)
			return
		}
	}
	return
}
func NuSerializeMap(v Map) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	for key, e := range v {
		tmp[key], err = nuconvWrapErr(nuplugin.ToValue(e))
		if err != nil {
			err = fmt.Errorf("serialize map entry (%s:%v): %w", key, e, err)
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}
`

	hashmap := MapType{
		ID: TypeID{
			Pkg:  "testpkg",
			Name: "Map",
		},
		Value: BuiltinTypeID("int"),
	}

	testType(t, "map.go", hashmap, expected)
}
