package logic

import (
	"github.com/dave/jennifer/jen"
)

// Primitive maps primitive nu <-> primitive go type
type Primitive struct {
	ID     TypeID
	GoType string
}

// ID returns the ID of the type
func (primitive Primitive) TypeID() TypeID {
	return primitive.ID
}

// Definition returns a statement which declares the type definition
func (primitive Primitive) Definition(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (primitive Primitive) Parser(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (primitive Primitive) Serializer(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}
