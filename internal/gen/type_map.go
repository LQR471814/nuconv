package gen

import (
	"github.com/dave/jennifer/jen"
)

type MapType struct {
	ID    TypeID
	Value TypeID
}

// ID returns the ID of the type
func (maptype MapType) TypeID() TypeID {
	return maptype.ID
}

// GoType returns the Go definition of this type
func (maptype MapType) GoType(prev *jen.Statement) *jen.Statement {
	return maptype.Value.Qual(prev.Map(jen.String()))
}

// Definition returns a statement which declares the type definition
//
// NOTE: the reason why it returns an entire statement instead of just the
// expression is so that if any additional variables necessary for the
// definition must be declared, they can be declared
func (maptype MapType) Definition(g GenContext) (out *jen.Statement, err error) {
	out = defTemplate(g, maptype).
		Qual(nu_types_pkg, "Record").
		Call(jen.Qual(nu_types_pkg, "RecordDef").Values())
	return
}

// Parser returns a statement which declares the parsing function
func (maptype MapType) Parser(g GenContext) (out *jen.Statement, err error) {
	cast, errHandle := castValueBlock(jen.Map(jen.String()).Qual(nu_pkg, "Value"))

	out = parserTypeTemplate(
		g, maptype,
		cast,
		errHandle,
		jen.Id(id_out).Op("=").Make(maptype.Value.Qual(jen.Map(jen.String()))),
		jen.For(jen.List(jen.Id(id_key), jen.Id(id_el)).Op(":=").Range().Id(id_typed)).Block(
			maptype.Value.ParserQual(
				jen.List(
					jen.Id(id_out).Index(jen.Id(id_key)),
					jen.Id(id_err_val),
				).Op("="),
				jen.Id(id_el),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit("parse map entry (%s:%v): %w"),
					jen.Id(id_key),
					jen.Id(id_el),
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
func (maptype MapType) Serializer(g GenContext) (out *jen.Statement, err error) {
	out = serializerTypeTemplate(
		g, maptype,
		jen.Id(id_tmp).Op(":=").Make(jen.Map(jen.String()).Qual(nu_pkg, "Value")),
		jen.For(jen.List(jen.Id(id_key), jen.Id(id_el)).Op(":=").Range().Id(id_value)).Block(
			maptype.Value.SerializerQual(
				jen.List(
					jen.Id(id_tmp).Index(jen.Id(id_key)),
					jen.Id(id_err_val),
				).Op("="),
				jen.Id(id_el),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit("serialize map entry (%s:%v): %w"),
					jen.Id(id_key),
					jen.Id(id_el),
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
