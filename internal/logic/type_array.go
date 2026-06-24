package logic

import "github.com/dave/jennifer/jen"

// ArrayType maps a go array <-> nu list
type ArrayType struct {
	ID          TypeID
	ElementType TypeID
	Length      int64
}

// ID returns the ID of the type
func (arraytype ArrayType) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (arraytype ArrayType) Definition() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (arraytype ArrayType) Parser() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (arraytype ArrayType) Serializer() *jen.Statement {
	panic("not implemented") // TODO: Implement
}
