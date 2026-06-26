package gen

import (
	"testing"
)

func TestRecordType(t *testing.T) {
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
	testType(t, "record.go", typ, "")
}
