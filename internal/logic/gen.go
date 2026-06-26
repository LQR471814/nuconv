package logic

import "github.com/dave/jennifer/jen"

const (
	nu_pkg       = "github.com/ainvaltin/nu-plugin"
	nu_types_pkg = "github.com/ainvaltin/nu-plugin/types"
)

const (
	id_i           = "i"
	id_el          = "e"
	id_out         = "out"
	id_value       = "v"
	id_ok          = "ok"
	id_typed       = "typed"
	id_err_val     = "err"
	id_tmp         = "tmp"
	id_wrap_err_fn = "nuconvWrapErr"
)

// v nu.Value
func nuValueParam(name string) *jen.Statement {
	return jen.Id(name).Qual(nu_pkg, "Value")
}

// err error
func errParam() *jen.Statement {
	return jen.Id(id_err_val).Id("error")
}

// typed, ok := v.(pkg.Type)
// if !ok { err = fmt.Errorf(...); return }
func parserCastValue(typ *jen.Statement) (cast, handle *jen.Statement) {
	cast = jen.List(jen.Id(id_typed), jen.Id(id_ok)).
		Op(":=").Id(id_value).Dot("Value").Assert(typ)
	handle = jen.If(jen.Op("!").Id(id_ok)).Block(
		jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
			jen.Lit("expected %T got %T"),
			jen.Id(id_typed),
			jen.Id(id_value),
		),
		jen.Return(),
	)
	return
}

// var NuDef... =
func typeDefTemplate(g GenContext, t Type) *jen.Statement {
	return g.Out.Var().Id(NewTypeDefID(t.TypeID())).Op("=")
}

//	func NuParser...(v nu.Value) (out ..., err error) {
//	  // cast type
//	}
func parserTypeTemplate(
	g GenContext,
	t Type,
	statements ...jen.Code,
) *jen.Statement {
	cast, handle := parserCastValue(jen.Index().Qual(nu_pkg, "Value"))

	block := []jen.Code{cast, handle}
	block = append(block, statements...)

	return g.Out.Func().
		Id(NewParserID(t.TypeID())).
		Params(nuValueParam(id_value)).
		Params(t.TypeID().Qual(jen.Id(id_out)), errParam()).
		Block(block...)
}

// func NuSerializer...(v nu.Value) (out ..., err error) {
// }
func serializerTypeTemplate(
	g GenContext,
	t Type,
	statements ...jen.Code,
) *jen.Statement {
	return g.Out.Func().
		Id(NewSerializerID(t.TypeID())).
		Params(t.TypeID().Qual(jen.Id(id_value))).
		Params(nuValueParam(id_out), errParam()).
		Block(statements...)
}
