package gen

import "github.com/dave/jennifer/jen"

const (
	nu_pkg       = "github.com/ainvaltin/nu-plugin"
	nu_types_pkg = "github.com/ainvaltin/nu-plugin/types"
)

const (
	id_T           = "T"
	id_i           = "i"
	id_el          = "e"
	id_out         = "out"
	id_value       = "v"
	id_ok          = "ok"
	id_typed       = "typed"
	id_err_val     = "err"
	id_tmp         = "tmp"
	id_wrap_err_fn = "nuconvWrapErr"
	id_try_cast_fn = "tryCast"
)

// v nu.Value
func nuValueParam(name string) *jen.Statement {
	return jen.Id(name).Qual(nu_pkg, "Value")
}

// err error
func errParam() *jen.Statement {
	return jen.Id(id_err_val).Id("error")
}

// typed, err := tryCast[pkg.Type](v.Value)
// if err != nil { err = fmt.Errorf("cast: %w", err); return }
func castValueBlock(typ *jen.Statement) (cast, errHandle *jen.Statement) {
	cast = jen.List(jen.Id(id_typed), jen.Id(id_err_val)).
		Op(":=").Id(id_try_cast_fn).Types(typ).Call(jen.Id(id_value).Dot("Value"))
	errHandle = jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
		jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
			jen.Lit("cast: %w"),
			jen.Id(id_err_val),
		),
		jen.Return(),
	)
	return
}

// var NuDef... =
func defTemplate(g GenContext, t Type) *jen.Statement {
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
	return g.Out.Func().
		Id(NewParserID(t.TypeID())).
		Params(nuValueParam(id_value)).
		Params(t.TypeID().Qual(jen.Id(id_out)), errParam()).
		Block(statements...)
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

func RenderHelpers(g GenContext) {
	wrapErrFn(g)
	tryCastFn(g)
}

func wrapErrFn(g GenContext) {
	g.Out.Func().
		Id(id_wrap_err_fn).
		Types(jen.Id(id_T).Any()).
		Params(jen.Id(id_value).Id(id_T)).
		Params(jen.Id(id_T), jen.Error()).
		Block(jen.Return(jen.Id(id_value), jen.Nil()))
}

func tryCastFn(g GenContext) {
	g.Out.Func().
		Id(id_try_cast_fn).
		Types(jen.Id(id_T).Any()).
		Params(jen.Id(id_value).Any()).
		Params(jen.Id(id_out).Id(id_T), jen.Id(id_err_val).Error()).
		Block(
			jen.List(jen.Id(id_out), jen.Id(id_ok)).
				Op(":=").
				Id(id_value).Assert(jen.Id(id_T)),
			jen.If(jen.Op("!").Id(id_ok)).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit("expected %T got %T"),
					jen.Id(id_out),
					jen.Id(id_value),
				),
				jen.Return(),
			),
			jen.Return(),
		)
}
