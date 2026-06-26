package logic

import "github.com/dave/jennifer/jen"

// OneofType maps either:
//
//   - *<type> -> oneof<type, nothing>
//   - struct satisfying "oneof" intf -> oneof<...>
type OneofType struct {
	ID   TypeID
	Alts []TypeID
}

// ID returns the ID of the type
func (oneoftype OneofType) TypeID() TypeID {
	return oneoftype.ID
}

// Definition returns a statement which declares the type definition
func (oneoftype OneofType) Definition(g GenContext) (out *jen.Statement, err error) {
	params := make([]jen.Code, len(oneoftype.Alts))
	for i, alt := range oneoftype.Alts {
		params[i] = alt.DefQual(nil)
	}
	out = defTemplate(g, oneoftype).
		Qual(nu_pkg, "Oneof").
		Call(params...)
	return
}

// Parser returns a statement which declares the parsing function
func (oneoftype OneofType) Parser(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}

// Serializer returns a statement which declares the serializer function
func (oneoftype OneofType) Serializer(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}
