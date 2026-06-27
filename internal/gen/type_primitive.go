package gen

import "github.com/dave/jennifer/jen"

type PrimitiveType struct {
	ID         TypeID
	GoTypeName string
}

// ID returns the ID of the type
func (primitive PrimitiveType) TypeID() TypeID {
	return primitive.ID
}

// GoType returns the Go definition of this type
func (primitive PrimitiveType) GoType(prev *jen.Statement) *jen.Statement {
	if prev == nil {
		prev = newStatement()
	}
	return prev.Id(primitive.GoTypeName)
}

// Definition returns a statement which declares the type definition
//
// NOTE: the reason why it returns an entire statement instead of just the
// expression is so that if any additional variables necessary for the
// definition must be declared, they can be declared
func (primitive PrimitiveType) Definition(g GenContext) (out *jen.Statement, err error) {
	out = primitive.ID.DefQual(nil)
	return
}

// Parser returns a statement which declares the parsing function
func (primitive PrimitiveType) Parser(g GenContext) (out *jen.Statement, err error) {
	cast, errHandle := castValueBlock(primitive.GoType(nil))
	out = parserTypeTemplate(g, primitive, cast, errHandle, jen.Return())
	return
}

// Serializer returns a statement which declares the serializer function
func (primitive PrimitiveType) Serializer(g GenContext) (out *jen.Statement, err error) {
	out = parserTypeTemplate(
		g, primitive,
		jen.Id(id_out).Op("=").Qual(nu_pkg, "ToValue").Call(jen.Id(id_value)),
		jen.Return(),
	)
	return
}
