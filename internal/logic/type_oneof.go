package logic

import "github.com/dave/jennifer/jen"

// OneofType maps either:
//
//   - *<type> -> oneof<type, nothing>
//   - struct satisfying "oneof" intf -> oneof<...>
type OneofType struct {
	ID           TypeID
	Alternatives []TypeID
}

// ID returns the ID of the type
func (oneoftype OneofType) TypeID() TypeID {
	panic("not implemented") // TODO: Implement
}

// Definition returns a statement which declares the type definition
func (oneoftype OneofType) Definition() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Parser returns a statement which declares the parsing function
func (oneoftype OneofType) Parser() *jen.Statement {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (oneoftype OneofType) Serializer() *jen.Statement {
	panic("not implemented") // TODO: Implement
}
