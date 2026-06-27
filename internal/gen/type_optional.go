package gen

import "github.com/dave/jennifer/jen"

type OptionalType struct {
	ID   TypeID
	Elem TypeID
}

// ID returns the ID of the type
func (optionaltype OptionalType) TypeID() TypeID {
	return optionaltype.ID
}

// Definition returns a statement which declares the type definition
//
// NOTE: the reason why it returns an entire statement instead of just the
// expression is so that if any additional variables necessary for the
// definition must be declared, they can be declared
func (optionaltype OptionalType) Definition(g GenContext) (out *jen.Statement, err error) {
	// oneof<elem, nothing>
	out = defTemplate(g, optionaltype).
		Qual(nu_types_pkg, "OneOf").
		Call(
			jen.Qual(nu_types_pkg, "Nothing").Call(),
			optionaltype.Elem.DefQual(nil),
		)
	return
}

// Parser returns a statement which declares the parsing function
func (optionaltype OptionalType) Parser(g GenContext) (out *jen.Statement, err error) {
	out = parserTypeTemplate(
		g, optionaltype,
		jen.If(jen.Id(id_value).Dot("Value").Op("==").Nil()).Block(
			jen.Return(),
		),
		optionaltype.Elem.Qual(jen.Var().Id(id_tmp)),
		optionaltype.Elem.ParserQual(
			jen.List(jen.Id(id_tmp), jen.Id(id_err_val)).Op("="),
			jen.Id(id_value).Dot("Value"),
		),
		jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
			jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
				jen.Lit("parse underlying: %w"),
				jen.Id(id_err_val),
			),
			jen.Return(),
		),
		jen.Id(id_out).Op("=").Op("&").Id(id_tmp),
		jen.Return(),
	)
	return
}

// Serializer returns a statement which declares the serializer function
func (optionaltype OptionalType) Serializer(g GenContext) (out *jen.Statement, err error) {
	out = serializerTypeTemplate(
		g, optionaltype,
		jen.If(jen.Id(id_value).Op("==").Nil()).Block(
			jen.Id(id_out).Op("=").Qual(nu_pkg, "ToValue").Call(jen.Nil()),
			jen.Return(),
		),
		optionaltype.Elem.SerializerQual(
			jen.List(jen.Id(id_out), jen.Id(id_err_val)).Op("="),
			jen.Op("*").Id(id_value),
		),
		jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
			jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
				jen.Lit("serialize underlying: %w"),
				jen.Id(id_err_val),
			),
			jen.Return(),
		),
		jen.Return(),
	)
	return
}
