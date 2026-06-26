package gen

import "github.com/dave/jennifer/jen"

// ArrayType maps a go array <-> nu list
type ArrayType struct {
	ID          TypeID
	ElementType TypeID
	Length      int64
}

// ID returns the ID of the type
func (arraytype ArrayType) TypeID() TypeID {
	return arraytype.ID
}

// Definition returns a statement which declares the type definition
func (arraytype ArrayType) Definition(g GenContext) (out *jen.Statement, err error) {
	out = defTemplate(g, arraytype).
		Qual(nu_types_pkg, "List").
		Call(arraytype.ElementType.DefQual(nil))
	return
}

// Parser returns a statement which declares the parsing function
func (arraytype ArrayType) Parser(g GenContext) (out *jen.Statement, err error) {
	arrayLen := int(arraytype.Length)
	out = parserTypeTemplate(
		g, arraytype, jen.Index().Qual(nu_pkg, "Value"),
		jen.If(
			jen.Len(jen.Id(id_typed)).
				Op("!=").
				Lit(arrayLen),
		).Block(
			jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
				jen.Lit("array has wrong number of elements, got (%d), expected (%d)"),
				jen.Id("len").Call(jen.Id(id_typed)),
				jen.Lit(arrayLen),
			),
			jen.Return(),
		),
		jen.Id(id_out).Op("=").Make(arraytype.ID.Qual(nil)),
		jen.For(jen.Id(id_i).Op(":=").Range().Lit(arrayLen)).Block(
			arraytype.ElementType.ParserQual(jen.List(
				jen.Id(id_out).Index(jen.Id(id_i)),
				jen.Id(id_err_val),
			).Op("="), jen.Id(id_typed).Index(jen.Id(id_i))),
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
func (arraytype ArrayType) Serializer(g GenContext) (out *jen.Statement, err error) {
	arrayLen := int(arraytype.Length)
	out = serializerTypeTemplate(
		g, arraytype,
		jen.Id(id_tmp).Op(":=").Make(
			jen.Index().Qual(nu_pkg, "Value"),
			jen.Lit(arrayLen),
		),
		jen.For(jen.Id(id_i).Op(":=").Range().Lit(arrayLen)).Block(
			arraytype.ElementType.SerializerQual(jen.List(
				jen.Id(id_tmp).Index(jen.Id(id_i)),
				jen.Id(id_err_val),
			).Op("="), jen.Id(id_value).Index(jen.Id(id_i))),
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
