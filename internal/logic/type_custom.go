package logic

import (
	"github.com/dave/jennifer/jen"
)

// CustomType maps nu-plugin go custom type <-> nu custom tpye
type CustomType struct {
	ID   TypeID
	Name string
	// parser and serializer for this type would only work if ToBaseValue does not error
}

// ID returns the ID of the type
func (customtype CustomType) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (customtype CustomType) Definition(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (customtype CustomType) Parser(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (customtype CustomType) Serializer(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

