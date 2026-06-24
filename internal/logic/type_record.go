package logic

import (
	"github.com/dave/jennifer/jen"
)

// RecordType maps go struct <-> nu record
type RecordType struct {
	ID     TypeID
	Fields []Field
}

// ID returns the ID of the type
func (recordtype RecordType) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (recordtype RecordType) Definition() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (recordtype RecordType) Parser() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (recordtype RecordType) Serializer() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Field maps a go struct field <-> nu field
type Field struct {
	// NuName is the name of the nushell-facing field name.
	NuName string
	// GoName is the golang-facing field name.
	GoName string
	// Type is the TypeID of the field's type.
	Type TypeID
}
