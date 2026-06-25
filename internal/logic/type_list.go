package logic

import (
	"github.com/dave/jennifer/jen"
)

// ListType maps go slice <-> nu list
type ListType struct {
	ID          TypeID
	ElementType TypeID
}

// ID returns the ID of the type
func (listtype ListType) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (listtype ListType) Definition(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (listtype ListType) Parser(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (listtype ListType) Serializer(g Generator) *jen.Statement {
	panic("not implemented") // TODO: Implement
}

