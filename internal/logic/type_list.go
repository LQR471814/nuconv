package logic

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
	out = g.Out.Var().
		Id(NewTypeDefID(listtype.ID)).
		Op("=").
		Qual(nu_types_pkg, "List").Call(listtype.ElementType.Qual(nil))
	return
}

// Parser returns a statement which declares the parsing function
func (listtype ListType) Parser(g GenContext) (out *jen.Statement, err error) {
	cast, handle := parserCastValue(jen.Index().Qual(nu_pkg, "Value"))

	out = g.Out.Func().
		Id(NewParserID(listtype.ID)).
		Params(nuValueParam(id_value)).
		Params(
			listtype.ID.Qual(jen.Id(id_out)),
			errParam(),
		).
		Block(
			cast,
			handle,
			jen.Id(id_out).Op("=").Make(listtype.ID.Qual(nil)),
			jen.For(jen.Id(id_i).Op(":=").Range().Len(jen.Id(id_typed))).Block(
				listtype.ElementType.ParserQual(jen.List(
					jen.Id(id_out).Index(jen.Id(id_i)),
					jen.Id(id_err_val),
				).Op("="), jen.Id(id_typed).Index(jen.Id(id_i))),
				jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
					jen.Return(),
				),
			),
			jen.Return(),
		)
	return
}

// Serializer returns a statement which declares the serializer function
func (listtype ListType) Serializer(g GenContext) (out *jen.Statement, err error) {
	panic("not implemented") // TODO: Implement
}
