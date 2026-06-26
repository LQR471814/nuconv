package gen

import (
	"testing"
)

func TestRecordType(t *testing.T) {
	const expected = `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDefRecord = types.Record(types.RecordDef{
	"age":  types.Int(),
	"name": types.String(),
})

func NuParseRecord(v nuplugin.Value) (out Record, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out.Name, err = tryCast[string](typed["name"])
	if err != nil {
		err = fmt.Errorf("parse Name (name): %w", err)
		return
	}
	out.Age, err = tryCast[int](typed["age"])
	if err != nil {
		err = fmt.Errorf("parse Age (age): %w", err)
		return
	}
	return
}
func NuSerializeRecord(v Record) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	tmp["name"], err = nuconvWrapErr(nuplugin.ToValue(v.Name))
	if err != nil {
		err = fmt.Errorf("serialize Name (name): %w", err)
		return
	}
	tmp["age"], err = nuconvWrapErr(nuplugin.ToValue(v.Age))
	if err != nil {
		err = fmt.Errorf("serialize Age (age): %w", err)
		return
	}
	out = nuplugin.ToValue(tmp)
	return
}
`
	typ := RecordType{
		ID: TypeID{
			Pkg:  "testpkg",
			Name: "Record",
		},
		Fields: []Field{
			{
				NuName: "name",
				GoName: "Name",
				Type:   BuiltinTypeID("string"),
			},
			{
				NuName: "age",
				GoName: "Age",
				Type:   BuiltinTypeID("int"),
			},
		},
	}
	testType(t, "record.go", typ, expected)
}
