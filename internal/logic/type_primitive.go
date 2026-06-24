package logic

import (
	"github.com/dave/jennifer/jen"
)

// Primitive maps primitive nu <-> primitive go type
type Primitive struct {
	ID     TypeID
	NuType string
	GoType string
}

// ID returns the ID of the type
func (primitive Primitive) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (primitive Primitive) Definition() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (primitive Primitive) Parser() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (primitive Primitive) Serializer() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

