package gen

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
	return listtype.ID
}

// Definition returns a statement which declares the type definition
func (listtype ListType) Definition(g GenContext) (out *jen.Statement, err error) {
	out = defTemplate(g, listtype).
		Qual(nu_types_pkg, "List").
		Call(listtype.ElementType.DefQual(nil))
	return
}

// Parser returns a statement which declares the parsing function
func (listtype ListType) Parser(g GenContext) (out *jen.Statement, err error) {
	cast, errHandle := castValueBlock(jen.Index().Qual(nu_pkg, "Value"))

	out = parserTypeTemplate(
		g, listtype,
		cast,
		errHandle,
		jen.Id(id_out).Op("=").Make(listtype.ID.Qual(nil)),
		jen.For(jen.List(jen.Id(id_i), jen.Id(id_el)).Op(":=").Range().Id(id_typed)).Block(
			listtype.ElementType.ParserQual(
				jen.List(jen.Id(id_out).Index(jen.Id(id_i)), jen.Id(id_err_val)).Op("="),
				jen.Id(id_el).Dot("Value"),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit("parse element (%d): %w"),
					jen.Id(id_i),
					jen.Id(id_err_val),
				),
				jen.Return(),
			),
		),
		jen.Return(),
	)
	return
}

// Serializer returns a statement which declares the serializer function
func (listtype ListType) Serializer(g GenContext) (out *jen.Statement, err error) {
	out = serializerTypeTemplate(
		g, listtype,
		jen.Id(id_tmp).Op(":=").Make(
			jen.Index().Qual(nu_pkg, "Value"),
			jen.Len(jen.Id(id_value)),
		),
		jen.For(jen.List(jen.Id(id_i), jen.Id(id_el)).Op(":=").Range().Id(id_value)).Block(
			listtype.ElementType.SerializerQual(
				jen.List(jen.Id(id_tmp).Index(jen.Id(id_i)), jen.Id(id_err_val)).Op("="),
				jen.Id(id_el),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit("serialize element (%d): %w"),
					jen.Id(id_i),
					jen.Id(id_err_val),
				),
				jen.Return(),
			),
		),
		jen.Id(id_out).Op("=").Qual(nu_pkg, "ToValue").Call(jen.Id(id_tmp)),
		jen.Return(),
	)
	return
}
